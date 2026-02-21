package worker

import (
	"context"
	"fmt"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/N3moAhead/bombahead/match_runner/internal/config"
	"github.com/N3moAhead/bombahead/match_runner/internal/match"
	"github.com/N3moAhead/bombahead/match_runner/internal/mq"
	"github.com/N3moAhead/bombahead/match_runner/internal/runner"
	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var log = logger.New("[Worker]")

const retryCountHeader = "x-match-retry-count"

// Worker processes matches from the queue
type Worker struct {
	config *config.Config
	mq     *mq.Client
	runner *runner.Runner
}

// New creates a new worker instance.
func New(cfg *config.Config) (*Worker, error) {
	mqClient, err := mq.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create MQ client: %w", err)
	}

	return &Worker{
		config: cfg,
		mq:     mqClient,
		runner: runner.New(),
	}, nil
}

// Run starts the worker's main loop.
// It blocks until a shutdown signal is received.
func (w *Worker) Run() (runErr error) {
	defer func() {
		if err := w.mq.Close(); err != nil {
			if runErr == nil {
				runErr = fmt.Errorf("worker stopped and failed to close mq client cleanly: %w", err)
				return
			}
			log.Error("Failed to close MQ client during shutdown: %v", err)
		}
	}()

	msgs, err := w.mq.ConsumeMatchMessages()
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info("Worker is waiting for matches. Press CTRL+C to exit.")

	for {
		select {
		case <-ctx.Done():
			log.Info("Worker is shutting down.")
			return nil
		case amqpErr := <-w.mq.ChannelClose():
			if amqpErr != nil {
				return fmt.Errorf("rabbitmq channel closed unexpectedly: %w", amqpErr)
			}
			return fmt.Errorf("rabbitmq channel closed")
		case amqpErr := <-w.mq.ConnectionClose():
			if amqpErr != nil {
				return fmt.Errorf("rabbitmq connection closed unexpectedly: %w", amqpErr)
			}
			return fmt.Errorf("rabbitmq connection closed")
		case msg, ok := <-msgs:
			if !ok {
				if ctx.Err() != nil {
					log.Info("Message channel closed during shutdown.")
					return nil
				}
				return fmt.Errorf("message channel closed by broker")
			}

			if err := w.handleMessage(ctx, msg); err != nil {
				return fmt.Errorf("failed while handling message: %w", err)
			}
		}
	}
}

// handleMessage decodes and processes a single match message from the queue.
func (w *Worker) handleMessage(ctx context.Context, msg amqp.Delivery) error {
	log.Info("Received a new match request.")

	var details match.Details
	if err := details.FromJSON(msg.Body); err != nil {
		log.Error("Failed to decode match details JSON: %v", err)
		return w.handleFailure(ctx, msg, nil, "decode_failed", err, false)
	}

	// Create a new context for this specific match.
	matchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := w.runner.RunMatch(matchCtx, &details, w.config.MatchHistoryDir)
	if err != nil {
		log.Error("Failed to run match '%s': %v", details.MatchID, err)
		return w.handleFailure(ctx, msg, &details, "match_run_failed", err, true)
	}

	log.Success("Successfully processed match '%s'.", result.MatchID)

	resultJSON, err := result.ToJSON()
	if err != nil {
		log.Error("Failed to encode match result: %v", err)
		return w.handleFailure(ctx, msg, &details, "result_encoding_failed", err, false)
	}

	if err := w.mq.PublishResultMessage(ctx, resultJSON); err != nil {
		log.Error("Failed to publish match result: %v", err)
		return w.handleFailure(ctx, msg, &details, "result_publish_failed", err, true)
	}

	log.Info("Published result for match '%s'.", result.MatchID)

	// Everything worked, acknowledge the message to delete it from the queue.
	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("failed to ack message for match '%s': %w", result.MatchID, err)
	}

	return nil
}

func (w *Worker) handleFailure(ctx context.Context, msg amqp.Delivery, details *match.Details, reason string, cause error, retryable bool) error {
	retryCount := readRetryCount(msg.Headers)
	matchID := "<unknown>"
	if details != nil && details.MatchID != "" {
		matchID = details.MatchID
	}

	if retryable && retryCount < w.config.MaxMatchRetries {
		nextRetryCount := retryCount + 1
		headers := cloneHeaders(msg.Headers)
		headers[retryCountHeader] = int32(nextRetryCount)

		if err := w.mq.PublishMatchMessageWithHeaders(ctx, msg.Body, headers); err != nil {
			if nackErr := nackMessage(msg, true); nackErr != nil {
				return fmt.Errorf("failed to republish retry for match '%s': %w (additionally failed to requeue original message: %v)", matchID, err, nackErr)
			}
			return nil
		}

		log.Warn(
			"Retrying match '%s' due to '%s' (%d/%d)",
			matchID,
			reason,
			nextRetryCount,
			w.config.MaxMatchRetries,
		)

		if err := msg.Ack(false); err != nil {
			return fmt.Errorf("failed to ack original message after scheduling retry for match '%s': %w", matchID, err)
		}
		return nil
	}

	failureEvent := &match.Failure{
		MatchID:    matchID,
		Reason:     reason,
		Error:      cause.Error(),
		RetryCount: retryCount,
		FailedAt:   time.Now().UTC(),
		Payload:    append([]byte(nil), msg.Body...),
	}

	failureJSON, err := failureEvent.ToJSON()
	if err != nil {
		if nackErr := nackMessage(msg, true); nackErr != nil {
			return fmt.Errorf("failed to encode failure event for match '%s': %w (additionally failed to requeue original message: %v)", matchID, err, nackErr)
		}
		return nil
	}

	if err := w.mq.PublishFailureMessage(ctx, failureJSON); err != nil {
		if nackErr := nackMessage(msg, true); nackErr != nil {
			return fmt.Errorf("failed to publish failure event for match '%s': %w (additionally failed to requeue original message: %v)", matchID, err, nackErr)
		}
		return nil
	}

	log.Error(
		"Match '%s' failed permanently after %d retries. Reason: %s",
		matchID,
		retryCount,
		reason,
	)

	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("failed to ack failed message for match '%s': %w", matchID, err)
	}

	return nil
}

func readRetryCount(headers amqp.Table) int {
	if headers == nil {
		return 0
	}

	raw, ok := headers[retryCountHeader]
	if !ok {
		return 0
	}

	switch v := raw.(type) {
	case int:
		if v < 0 {
			return 0
		}
		return v
	case int8:
		if v < 0 {
			return 0
		}
		return int(v)
	case int16:
		if v < 0 {
			return 0
		}
		return int(v)
	case int32:
		if v < 0 {
			return 0
		}
		return int(v)
	case int64:
		if v < 0 {
			return 0
		}
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || parsed < 0 {
			return 0
		}
		return parsed
	default:
		return 0
	}
}

func cloneHeaders(headers amqp.Table) amqp.Table {
	cloned := amqp.Table{}
	for key, value := range headers {
		cloned[key] = value
	}
	return cloned
}

func nackMessage(msg amqp.Delivery, requeue bool) error {
	if err := msg.Nack(false, requeue); err == nil {
		return nil
	}

	if err := msg.Reject(requeue); err == nil {
		return nil
	}

	return fmt.Errorf("nack and reject both failed")
}
