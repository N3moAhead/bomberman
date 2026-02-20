package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/N3moAhead/bombahead/match_runner/internal/config"
	"github.com/N3moAhead/bombahead/match_runner/internal/match"
	"github.com/N3moAhead/bombahead/match_runner/internal/mq"
	"github.com/N3moAhead/bombahead/match_runner/pkg/logger"
	"github.com/google/uuid"
)

var log = logger.New("[Match_Initiator]")

func main() {
	serverImage := flag.String("server", "docker.io/nemoahead/bomberman-os-server:latest", "Server docker image")
	client1Image := flag.String("client1", "ghcr.io/n3moahead/bomber:self-destruct", "Client 1 docker image")
	client2Image := flag.String("client2", "ghcr.io/n3moahead/bomber:idle", "Client 2 docker image")
	flag.Parse()

	log.Info("Match Initiator is starting...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load configuration: %v", err))
	}

	mqClient, err := mq.NewClient(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create MQ client: %v", err))
	}
	defer mqClient.Close()

	matchID := uuid.New().String()
	details := match.Details{
		MatchID:      matchID,
		ServerImage:  *serverImage,
		Client1Image: *client1Image,
		Client2Image: *client2Image,
	}

	jsonData, err := details.ToJSON()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to marshal match details: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mqClient.PublishMatchMessage(ctx, jsonData); err != nil {
		log.Fatal(fmt.Sprintf("Failed to publish match message: %v", err))
	}

	log.Success("Successfully published match %s", matchID)
	log.Info("Server Image: %s", details.ServerImage)
	log.Info("Client 1 Image: %s", details.Client1Image)
	log.Info("Client 2 Image: %s", details.Client2Image)
}
