package cfg

import (
	"os"

	"github.com/N3moAhead/bomberman/website/pkg/logger"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var log = logger.New("[Config]")

type Config struct {
	DBURI              string
	Port               string
	GithubCLientId     string
	GithubClientSecret string
	GithubScopes       []string
	GithubEndpoint     oauth2.Endpoint
	NextAuthUrl        string
	SessionSecret      string
	IsProduction       bool
	CSRFAuthKey        string
	BaseURL            string
	RabbitMQURL        string
	MatchQueue         string
	ResultQueue        string
}

func Load() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("No .env file found...")
	}

	dbUri := os.Getenv("DBURI")
	if dbUri == "" {
		hasToBeSet("DBURI")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}

	rabbitMqUrl := os.Getenv("RABBITMQ_URL")
	if rabbitMqUrl == "" {
		log.Error("The env Variable RABBITMQ_URL has to be set!")
	}

	matchQueue := os.Getenv("RABBITMQ_MATCH_QUEUE")
	if matchQueue == "" {
		matchQueue = "bomberman.matches.pending"
	}

	resultQueue := os.Getenv("RABBITMQ_RESULT_QUEUE")
	if resultQueue == "" {
		resultQueue = "bomberman.matches.results"
	}

	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	if githubClientId == "" {
		hasToBeSet("GITHUB_CLIENT_ID")
	}

	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	if githubClientSecret == "" {
		hasToBeSet("GITHUB_CLIENT_SECRET")
	}

	nextAuthUrl := os.Getenv("NEXT_AUTH_URL")
	if nextAuthUrl == "" {
		hasToBeSet("NEXT_AUTH_URL")
	}

	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		hasToBeSet("SESSION_SECRET")
	}

	isProduction := os.Getenv("IS_PRODUCTION") == "true"

	csrfAuthKey := os.Getenv("CSRF_AUTH_KEY")
	if csrfAuthKey == "" {
		hasToBeSet("CSRF_AUTH_KEY")
	}

	return &Config{
		DBURI:              dbUri,
		Port:               port,
		GithubCLientId:     githubClientId,
		GithubClientSecret: githubClientSecret,
		GithubEndpoint:     github.Endpoint,
		NextAuthUrl:        nextAuthUrl + "/auth/github/callback",
		SessionSecret:      sessionSecret,
		IsProduction:       isProduction,
		CSRFAuthKey:        csrfAuthKey,
		BaseURL:            nextAuthUrl,
		RabbitMQURL:        rabbitMqUrl,
		MatchQueue:         matchQueue,
		ResultQueue:        resultQueue,
	}
}

func hasToBeSet(name string) {
	log.Fatal("The env-variable ", name, " has to be set!")
}
