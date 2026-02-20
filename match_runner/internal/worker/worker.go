package worker

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/N3moAhead/bombahead/match_runner/internal/config"
	"github.com/N3moAhead/bombahead/match_runner/internal/match"
	"github.com/N3moAhead/bombahead/match_runner/internal/mq"
	"github.com/N3moAhead/bombahead/match_runner/internal/runner"
	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var log = logger.New("[Worker]")

// Worker processes matches from the queue.
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
		if nackErr := nackMessage(msg, false); nackErr != nil {
			return fmt.Errorf("decode failed: %w (additionally failed to reject message: %v)", err, nackErr)
		}
		return nil
	}

	// Create a new context for this specific match.
	matchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := w.runner.RunMatch(matchCtx, &details)
	if err != nil {
		log.Error("Failed to run match '%s': %v", details.MatchID, err)
		// Reject the message. Don't requeue to avoid poison pills.
		if nackErr := nackMessage(msg, false); nackErr != nil {
			return fmt.Errorf("match run failed: %w (additionally failed to reject message: %v)", err, nackErr)
		}
		return nil
	}

	log.Success("Successfully processed match '%s'.", result.MatchID)

	resultJSON, err := result.ToJSON()
	if err != nil {
		log.Error("Failed to encode match result: %v", err)
		// The match ran, but we can't report it.
		if nackErr := nackMessage(msg, false); nackErr != nil {
			return fmt.Errorf("result encoding failed: %w (additionally failed to reject message: %v)", err, nackErr)
		}
		return nil
	}

	if err := w.mq.PublishResultMessage(ctx, resultJSON); err != nil {
		log.Error("Failed to publish match result: %v", err)
		// Requeue the message so we can try publishing the result again.
		if nackErr := nackMessage(msg, true); nackErr != nil {
			return fmt.Errorf("result publishing failed: %w (additionally failed to requeue message: %v)", err, nackErr)
		}
		return nil
	}

	log.Info("Published result for match '%s'.", result.MatchID)

	// Everything worked, acknowledge the message to delete it from the queue.
	if err := msg.Ack(false); err != nil {
		return fmt.Errorf("failed to ack message for match '%s': %w", result.MatchID, err)
	}

	return nil
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
