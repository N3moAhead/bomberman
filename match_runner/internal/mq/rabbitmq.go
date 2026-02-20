package mq

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/N3moAhead/bombahead/match_runner/internal/config"
	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var log = logger.New("[MQ]")

const (
	connectRetries = 5
	retryDelay     = 5 * time.Second
)

// Client manages the connection and channel to a RabbitMQ server
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  *config.Config

	connClose <-chan *amqp.Error
	chanClose <-chan *amqp.Error
}

// NewClient establishes a connection to RabbitMQ
// and prepares the channel and queues.
func NewClient(cfg *config.Config) (*Client, error) {
	log.Info("Connecting to RabbitMQ at %s", cfg.RabbitMQURL)

	var (
		conn    *amqp.Connection
		lastErr error
	)

	for attempt := 1; attempt <= connectRetries; attempt++ {
		newConn, err := amqp.Dial(cfg.RabbitMQURL)
		if err == nil {
			conn = newConn
			break
		}

		lastErr = err
		if attempt < connectRetries {
			log.Warn(
				"Failed to connect to RabbitMQ (attempt %d/%d). Retrying in %s: %v",
				attempt,
				connectRetries,
				retryDelay,
				err,
			)
			time.Sleep(retryDelay)
		}
	}

	if conn == nil {
		return nil, fmt.Errorf("could not establish a connection to RabbitMQ after %d attempts: %w", connectRetries, lastErr)
	}

	ch, err := conn.Channel()
	if err != nil {
		if closeErr := conn.Close(); closeErr != nil && !errors.Is(closeErr, amqp.ErrClosed) {
			return nil, fmt.Errorf("failed to open a channel: %w (additionally failed closing connection: %v)", err, closeErr)
		}
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare durable queues to ensure they survive broker restarts
	queues := []string{cfg.MatchQueue, cfg.ResultQueue}
	for _, q := range queues {
		if _, err = ch.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			closeErr := closeAMQP(ch, conn)
			if closeErr != nil {
				return nil, fmt.Errorf("failed to declare queue '%s': %w (additionally failed during cleanup: %v)", q, err, closeErr)
			}
			return nil, fmt.Errorf("failed to declare queue '%s': %w", q, err)
		}
		log.Info("Queue '%s' is ready.", q)
	}

	log.Success("RabbitMQ client initialized successfully")

	return &Client{
		conn:      conn,
		ch:        ch,
		cfg:       cfg,
		connClose: conn.NotifyClose(make(chan *amqp.Error, 1)),
		chanClose: ch.NotifyClose(make(chan *amqp.Error, 1)),
	}, nil
}

// ConsumeMatchMessages starts consuming messages from the match queue
// It returns a channel of deliveries for the caller to process
func (c *Client) ConsumeMatchMessages() (<-chan amqp.Delivery, error) {
	if c == nil || c.ch == nil {
		return nil, fmt.Errorf("rabbitmq client is not initialized")
	}

	// Set prefetch count to 1 to ensure this runner only takes one message at a time
	if err := c.ch.Qos(1, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	consumerTag := fmt.Sprintf("match-runner-%d-%d", os.Getpid(), time.Now().UnixNano())

	msgs, err := c.ch.Consume(
		c.cfg.MatchQueue,
		consumerTag, // consumer tag must be unique across concurrently running workers
		false,       // auto-ack (we will manually ack/nack)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Info("Registered consumer '%s' on queue '%s' with prefetch=1", consumerTag, c.cfg.MatchQueue)

	return msgs, nil
}

// ChannelClose returns channel close notifications from RabbitMQ.
func (c *Client) ChannelClose() <-chan *amqp.Error {
	return c.chanClose
}

// ConnectionClose returns connection close notifications from RabbitMQ.
func (c *Client) ConnectionClose() <-chan *amqp.Error {
	return c.connClose
}

// PublishResultMessage publishes the result of a match to the result queue.
func (c *Client) PublishResultMessage(ctx context.Context, body []byte) error {
	if c == nil || c.ch == nil {
		return fmt.Errorf("rabbitmq client is not initialized")
	}

	log.Info("Publishing result to queue '%s'", c.cfg.ResultQueue)
	return c.ch.PublishWithContext(ctx,
		"",                // exchange (default)
		c.cfg.ResultQueue, // routing key (queue name)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Make message persistent
			Body:         body,
		})
}

// PublishMatchMessage publishes a match request to the match queue.
func (c *Client) PublishMatchMessage(ctx context.Context, body []byte) error {
	if c == nil || c.ch == nil {
		return fmt.Errorf("rabbitmq client is not initialized")
	}

	log.Info("Publishing match request to queue '%s'", c.cfg.MatchQueue)
	return c.ch.PublishWithContext(ctx,
		"",               // exchange (default)
		c.cfg.MatchQueue, // routing key (queue name)
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Make message persistent
			Body:         body,
		})
}

// Close gracefully closes the channel and connection
func (c *Client) Close() error {
	if c == nil {
		return nil
	}

	if err := closeAMQP(c.ch, c.conn); err != nil {
		return fmt.Errorf("failed to close rabbitmq resources: %w", err)
	}

	log.Info("RabbitMQ connection closed.")
	return nil
}

func closeAMQP(ch *amqp.Channel, conn *amqp.Connection) error {
	var errs []error

	if ch != nil {
		if err := ch.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			errs = append(errs, fmt.Errorf("channel close: %w", err))
		}
	}

	if conn != nil {
		if err := conn.Close(); err != nil && !errors.Is(err, amqp.ErrClosed) {
			errs = append(errs, fmt.Errorf("connection close: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
