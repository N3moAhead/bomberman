package main

import (
	"fmt"

	"github.com/N3moAhead/bomberman/match_runner/internal/config"
	"github.com/N3moAhead/bomberman/match_runner/internal/worker"
	"github.com/N3moAhead/bomberman/match_runner/pkg/logger"
)

var log = logger.New("[Match_Runner]")

func main() {
	log.Info("Match Runner is starting up...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	// The worker contains the main application logic
	w, err := worker.New(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create a new worker: %v", err))
	}

	if err := w.Run(); err != nil {
		log.Fatal(fmt.Sprintf("Worker stopped with an error: %v", err))
	}

	log.Info("Match Runner has shut down gracefully.")
}
