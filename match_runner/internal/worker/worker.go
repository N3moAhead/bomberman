package worker

import (
	"context"
	"fmt"
	"os"
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

// Run starts the worker's main loop
// It blocks until a shutdown signal is received
func (w *Worker) Run() error {
	defer w.mq.Close()

	msgs, err := w.mq.ConsumeMatchMessages()
	if err != nil {
		return fmt.Errorf("failed to start consuming messages: %w", err)
	}

	// TODO i should add somehting so that endless matches
	// end after some time...
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Info("Shutdown signal received, stopping worker...")
		cancel()
	}()

	log.Info("Worker is waiting for matches. Press CTRL+C to exit.")

	for {
		select {
		case <-ctx.Done():
			log.Info("Worker is shutting down.")
			return nil
		case msg, ok := <-msgs:
			if !ok {
				log.Info("Message channel closed by broker. Shutting down.")
				return nil
			}
			w.handleMessage(ctx, msg)
		}
	}
}

// handleMessage decodes and processes a single match message from the queue
func (w *Worker) handleMessage(ctx context.Context, msg amqp.Delivery) {
	log.Info("Received a new match request.")

	var details match.Details
	if err := details.FromJSON(msg.Body); err != nil {
		log.Error("Failed to decode match details JSON: %v", err)
		err := msg.Nack(false, false)
		if err != nil {
			log.Errorln("Failed to nack msg ", err)
		}
		return
	}

	// Create a new context for this specific match
	matchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := w.runner.RunMatch(matchCtx, &details)
	if err != nil {
		log.Error("Failed to run match '%s': %v", details.MatchID, err)
		// Reject the message. Don't requeue to avoid poison pills
		err := msg.Nack(false, false)
		if err != nil {
			log.Errorln("Failed to nack msg ", err)
		}
		return
	}

	log.Success("Successfully processed match '%s'.", result.MatchID)

	resultJSON, err := result.ToJSON()
	if err != nil {
		log.Error("Failed to encode match result: %v", err)
		// This is a weird state. The match ran, but we can't report it
		// Nack and don't requeue
		err := msg.Nack(false, false)
		if err != nil {
			log.Errorln("Failed to nack msg ", err)
		}
		return
	}

	if err := w.mq.PublishResultMessage(ctx, resultJSON); err != nil {
		log.Error("Failed to publish match result: %v", err)
		// Requeue the message so we can try publishing the result again
		err := msg.Nack(false, true)
		if err != nil {
			log.Errorln("Failed to nack message ", err)
		}
		return
	}

	log.Info("Published result for match '%s'.", result.MatchID)

	// Everything worked so we can just
	// acknowledge the message to delete it from the queue
	err = msg.Ack(false)
	if err != nil {
		log.Errorln("Failed to Ack the msg ", err)
	}
}
