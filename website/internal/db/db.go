package db

import (
	"time"

	"github.com/N3moAhead/bombahead/website/internal/cfg"
	"github.com/N3moAhead/bombahead/website/internal/models"
	"github.com/N3moAhead/bombahead/website/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB
var log = logger.New("[DB]")

func Init(cfg *cfg.Config) {
	log.Infoln("Connecting to DB...")
	maxTries := 15
	for range maxTries {
		database, err := gorm.Open(postgres.Open(cfg.DBURI), &gorm.Config{})
		if err == nil {
			Conn = database
			log.Successln("Successfully connected to the db")
			break
		}
		log.Warnln("Could not establish a db connection. Trying again in 5 seconds.", err)
		time.Sleep(5 * time.Second)
	}
	if Conn == nil {
		log.Fatal("Could not establish a db connection. Shutting down")
	}

	log.Infoln("Auto-migrating models...")
	err := Conn.AutoMigrate(&models.User{}, &models.Bot{}, &models.Match{})
	if err != nil {
		log.Fatal("Could not auto-migrate models", err)
	}
	log.Successln("Successfully auto-migrated models")
}
