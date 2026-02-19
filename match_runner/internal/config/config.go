package config

import (
	"os"

	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
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
		log.Warn("No .env file found...")
	}

	// The loader is pretty simple for now because i don't need much
	// In the future i might extend those things with some functions
	// for loading config files ord using flags
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		log.Fatal("The env Variable RABBITMQ_URL has to be set!")
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
