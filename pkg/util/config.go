package util

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	db_host      string
	db_user      string
	db_password  string
	db_namespace string
	db_name      string

	open_ai_api_key     string
	discord_api_key     string
	dev_discord_api_key string
	finn_hub_api_key    string
	elevenlabs_api_key  string
	perplexity_api_key  string
)

func ReadEnv() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	envFilePath := filepath.Join(cwd, ".env")
	err = godotenv.Overload(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db_host = os.Getenv("DB_HOST")
	db_user = os.Getenv("DB_USER")
	db_password = os.Getenv("DB_PASSWORD")
	db_namespace = os.Getenv("DB_NAMESPACE")
	db_name = os.Getenv("DB_NAME")

	open_ai_api_key = os.Getenv("OPEN_AI_API_KEY")
	discord_api_key = os.Getenv("DISCORD_API_KEY")
	dev_discord_api_key = os.Getenv("DEV_DISCORD_API_KEY")
	finn_hub_api_key = os.Getenv("FINNHUB_API_KEY")
	elevenlabs_api_key = os.Getenv("ELEVEN_LABS_API_KEY")
	perplexity_api_key = os.Getenv("PERPLEXITY_API_KEY")
}

func GetDBHost() string {
	return db_host
}

func GetDBUser() string {
	return db_user
}

func GetDBPassword() string {
	return db_password
}

func GetDBNamespace() string {
	return db_namespace
}

func GetDBName() string {
	return db_name
}

func GetOpenAIKey() string {
	return open_ai_api_key
}

func GetPerplexityAPIKey() string {
	return perplexity_api_key
}

func GetDiscordKey() string {
	return discord_api_key
}

func GetDevDiscordKey() string {
	return dev_discord_api_key
}

func GetFinnHubKey() string {
	return finn_hub_api_key
}

func GetElevenLabsKey() string {
	return elevenlabs_api_key
}

func UserIsAdmin(userId string) bool {
	return false
}
