package mq

import (
	"context"
	"fmt"
	"time"

	"github.com/N3moAhead/bomberman/match_runner/internal/config"
	"github.com/N3moAhead/bomberman/match_runner/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

var log = logger.New("[MQ]")

// Client manages the connection and channel to a RabbitMQ server
type Client struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  *config.Config
}

// NewClient establishes a connection to RabbitMQ
// and prepares the channel and queues
func NewClient(cfg *config.Config) (*Client, error) {
	log.Info("Connecting to RabbitMQ at %s", cfg.RabbitMQURL)

	maxTries := 5

	var conn *amqp.Connection
	for range maxTries {
		newConn, err := amqp.Dial(cfg.RabbitMQURL)
		if err == nil {
			conn = newConn
			break
		}
		log.Warn("Failed to connect to RabbitMQ Trying again in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}
	if conn == nil {
		return nil, fmt.Errorf("Could not establish a connection to the rabbitMQ")
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare durable queues to ensure they survive broker restarts
	queues := []string{cfg.MatchQueue, cfg.ResultQueue}
	for _, q := range queues {
		_, err = ch.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, fmt.Errorf("failed to declare queue '%s': %w", q, err)
		}
		log.Info("Queue '%s' is ready.", q)
	}

	log.Success("RabbitMQ client initialized successfully")
	return &Client{conn: conn, ch: ch, cfg: cfg}, nil
}

// ConsumeMatchMessages starts consuming messages from the match queue
// It returns a channel of deliveries for the caller to process
func (c *Client) ConsumeMatchMessages() (<-chan amqp.Delivery, error) {
	// Set prefetch count to 1 to ensure this runner only takes one message at a time
	if err := c.ch.Qos(1, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return c.ch.Consume(
		c.cfg.MatchQueue,
		"match-runner", // consumer tag
		false,          // auto-ack (we will manually ack/nack)
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
}

// PublishResultMessage publishes the result of a match to the result queue
func (c *Client) PublishResultMessage(ctx context.Context, body []byte) error {
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

// PublishMatchMessage publishes a match request to the match queue
func (c *Client) PublishMatchMessage(ctx context.Context, body []byte) error {
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
func (c *Client) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
	log.Info("RabbitMQ connection closed.")
}
