package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/message"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/internal/mq"
	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/google/uuid"
	"github.com/intinig/go-openskill/rating"
	"github.com/intinig/go-openskill/types"
	"github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var log = logger.New("[Matchmaker]")

type UnmatchedBots struct {
	Bot1 uint
	Bot2 uint
}

func main() {
	config := cfg.Load()
	db := connectToDB(config)
	mqClient := connectToMQ(config)
	defer mqClient.Close()

	msgs, err := mqClient.ConsumeResultMessages()
	if err != nil {
		log.Errorln("Failed to start consuming messages: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown handling
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Info("Shutdown signal received, stopping match maker...")
		cancel()
	}()

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ctx.Done():
			log.Info("Matchmaker is shutting down!")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Info("Message channel closed by broker. Shutting down.")
				return
			}
			handleResultMessage(msg, db)
		case <-ticker.C:
			unmatchedBots := getUnmatchedBots(db)
			for _, match := range unmatchedBots {
				startNewMatch(match.Bot1, match.Bot2, mqClient, db)
			}
		}
	}
}

func startNewMatch(bot1ID uint, bot2ID uint, mqClient *mq.Client, db *gorm.DB) {
	var bot1 models.Bot
	err := db.Where("id = ?", bot1ID).First(&bot1).Error
	if err != nil {
		log.Error("Could not load bot with id %d! Stopping match start!", bot1ID)
		return
	}
	var bot2 models.Bot
	err = db.Where("id = ?", bot2ID).First(&bot2).Error
	if err != nil {
		log.Error("Could not load bot with id %d! Stopping match start!", bot2ID)
		return
	}

	matchID := uuid.New().String()
	details := message.Details{
		MatchID:      matchID,
		ServerImage:  "docker.io/nemoahead/bomberman-os-server:latest",
		Client1Image: bot1.DockerHubUrl,
		Client2Image: bot2.DockerHubUrl,
	}

	jsonData, err := details.ToJSON()
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to marshal match details: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mqClient.PublishMatchMessage(ctx, jsonData); err != nil {
		log.Errorln(fmt.Sprintf("Failed to publish match message: %v", err))
		return
	}

	// Saving the current state to the DB
	newMatch := &models.Match{
		MatchID: matchID,
		Bot1ID:  bot1ID,
		Bot2ID:  bot2ID,
		Status:  models.PENDING,
	}

	err = db.Create(newMatch).Error
	if err != nil {
		log.Error("Failed to save match to the DB")
		// TODO i really need a goood way to fix this issue! To make it way more stable....
		return
	}

	log.Successln("Successfully created new Match! YAY")
}

func handleResultMessage(msg amqp091.Delivery, db *gorm.DB) error {
	log.Info("Received a match result.")

	var matchResult message.Result
	err := json.Unmarshal(msg.Body, &matchResult)
	if err != nil {
		log.Error("Failed to process Match Results")
		msg.Nack(false, false)
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var dbMatch models.Match
		err = tx.Where("match_id = ?", matchResult.MatchID).Preload("Bot1").Preload("Bot2").First(&dbMatch).Error
		if err != nil {
			log.Error("Failed to find the corresponding db match to the received match result")
			msg.Nack(false, false)
			return err
		}

		// TODO
		// A probably pretty useless and stupid way to do it
		// All the time im just marshalling and unmarshalling the same JSON
		// So i could definitly improve it...
		// But diffrent things matter more at the moment so i will just leave it for the second...
		historyJson, err := json.Marshal(matchResult.Log)
		if err != nil {
			log.Error("Failed to marshal match history")
			msg.Nack(false, false)
			return err
		}
		dbMatch.History = historyJson
		dbMatch.Status = models.FINISHED

		var winner, loser models.Bot
		var options *types.OpenSkillOptions
		// Calculate Winner...
		switch matchResult.Winner {
		case dbMatch.Bot1.DockerHubUrl:
			// Bot1 Won
			winner, loser = dbMatch.Bot1, dbMatch.Bot2
			dbMatch.WinnerState = models.BOT1WIN
			options = nil
		case dbMatch.Bot2.DockerHubUrl:
			// Bot2 Won
			winner, loser = dbMatch.Bot1, dbMatch.Bot2
			dbMatch.WinnerState = models.BOT2WIN
			options = nil
		default:
			// Draw
			// winner or loser doesnt matter here
			winner, loser = dbMatch.Bot1, dbMatch.Bot2
			dbMatch.WinnerState = models.DRAW
			// both will just receive the same score
			options = &types.OpenSkillOptions{
				Score: []int{1, 1},
			}
		}

		tx.Save(&dbMatch)

		winnerRating := winner.ToRating()
		loserRating := loser.ToRating()

		teams := []types.Team{
			{winnerRating},
			{loserRating},
		}

		newRatings := rating.Rate(teams, options)

		winner.ApplyRating(newRatings[0][0])
		loser.ApplyRating(newRatings[1][0])

		if err := tx.Save(&winner).Error; err != nil {
			return err
		}
		if err := tx.Save(&loser).Error; err != nil {
			return err
		}

		msg.Ack(false)
		return nil
	})
}

func connectToMQ(config *cfg.Config) *mq.Client {
	maxTries := 5
	log.Infoln("Trying to connect to the RabbitMQ")
	for range maxTries {
		mqClient, err := mq.NewClient(config)
		if err == nil {
			log.Successln("Successfully conntected to the RabbitMQ")
			return mqClient
		}
		log.Warn("Could not connect to RabbitMQ trying again in 5s")
	}
	log.Fatal("Failed to conect to rabbitMQ")
	return nil
}

func connectToDB(config *cfg.Config) *gorm.DB {
	maxTries := 5
	log.Infoln("Establishing a connection to the db")
	for range maxTries {
		db, err := gorm.Open(postgres.Open(config.DBURI), &gorm.Config{})
		if err == nil {
			log.Successln("Successfully connected to the db")
			return db
		}
		log.Warn("Could not establish a db connection trying again in 5s")
		time.Sleep(5 * time.Second)
	}
	log.Fatal("Could not establish a database connection")
	return nil
}

func getUnmatchedBots(db *gorm.DB) []UnmatchedBots {
	var unmatchedBots []UnmatchedBots
	err := db.WithContext(context.Background()).Raw(`
			SELECT b1.id AS "Bot1", b2.id AS "Bot2"
			FROM bots b1
			JOIN bots b2 ON b1.id < b2.id

			EXCEPT

			SELECT
			    LEAST(bot1_id, bot2_id),
			    GREATEST(bot1_id, bot2_id)
			FROM matches;
		`).Scan(&unmatchedBots).Error

	if err != nil {
		log.Errorln(err)
	}

	return unmatchedBots
}
