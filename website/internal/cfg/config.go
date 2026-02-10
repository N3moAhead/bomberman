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

	return &Config{
		DBURI:              dbUri,
		Port:               port,
		GithubCLientId:     githubClientId,
		GithubClientSecret: githubClientSecret,
		GithubEndpoint:     github.Endpoint,
		NextAuthUrl:        nextAuthUrl + "/auth/github/callback",
	}
}

func hasToBeSet(name string) {
	log.Fatal("The env-variable", name, "has to be set!")
}
