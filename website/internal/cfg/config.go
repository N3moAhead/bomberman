package cfg

import (
	"os"

	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/joho/godotenv"
)

var log = logger.New("[Config]")

type Config struct {
	DBURI string
	Port  string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("No .env file found...")
	}

	dbUri := os.Getenv("DBURI")
	if dbUri == "" {
		log.Fatal("The env-variable DBURI has to be set!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	return &Config{DBURI: dbUri, Port: port}
}
