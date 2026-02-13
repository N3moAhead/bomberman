package db

import (
	"github.com/N3moAhead/bomberman/website/internal/cfg"
	"github.com/N3moAhead/bomberman/website/internal/models"
	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Conn *gorm.DB
var log = logger.New("[DB]")

func Init(cfg *cfg.Config) {
	log.Infoln("Connecting to DB...")
	database, err := gorm.Open(postgres.Open(cfg.DBURI), &gorm.Config{})
	if err != nil {
		log.Fatal("Could not establish a db connection", err)
	}
	Conn = database
	log.Successln("Successfully connected to the db")

	log.Infoln("Auto-migrating models...")
	err = Conn.AutoMigrate(&models.User{}, &models.Bot{})
	if err != nil {
		log.Fatal("Could not auto-migrate models", err)
	}
	log.Successln("Successfully auto-migrated models")
}
