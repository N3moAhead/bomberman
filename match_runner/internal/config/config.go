package config

import (
	"os"

	"github.com/N3moAhead/bomberman/match_runner/pkg/logger"
	"github.com/joho/godotenv"
)

var log = logger.New("[Config]")

// Config holds the configuration for the match runner
type Config struct {
	RabbitMQURL string
	MatchQueue  string
	ResultQueue string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Read the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// The loader is pretty simple for now because i don't need much
	// In the future i might extend those things with some functions
	// for loading config files ord using flags
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	matchQueue := os.Getenv("RABBITMQ_MATCH_QUEUE")
	if matchQueue == "" {
		matchQueue = "bomberman.matches.pending"
	}

	resultQueue := os.Getenv("RABBITMQ_RESULT_QUEUE")
	if resultQueue == "" {
		resultQueue = "bomberman.matches.results"
	}

	return &Config{
		RabbitMQURL: url,
		MatchQueue:  matchQueue,
		ResultQueue: resultQueue,
	}, nil
}
