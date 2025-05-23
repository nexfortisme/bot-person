package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	open_ai_api_key     string
	discord_api_key     string
	dev_discord_api_key string
	elevenlabs_api_key  string
	perplexity_api_key  string
	bot_open_ai_model   string
	open_ai_model       string
	admins              []string
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

	open_ai_api_key = os.Getenv("OPEN_AI_API_KEY")
	discord_api_key = os.Getenv("DISCORD_API_KEY")
	dev_discord_api_key = os.Getenv("DEV_DISCORD_API_KEY")
	elevenlabs_api_key = os.Getenv("ELEVEN_LABS_API_KEY")
	perplexity_api_key = os.Getenv("PERPLEXITY_API_KEY")
	bot_open_ai_model = os.Getenv("BOT_OPEN_AI_MODEL")
	open_ai_model = os.Getenv("OPEN_AI_MODEL")
	admins = strings.Split(os.Getenv("ADMINS"), ",")
}

func GetOpenAIKey() string {
	return open_ai_api_key
}

func GetOpenAIModel() string {
	return open_ai_model
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

func GetBotOpenAIModel() string {
	return bot_open_ai_model
}

func GetElevenLabsKey() string {
	return elevenlabs_api_key
}

func UserIsAdmin(userId string) bool {
	for _, admin := range admins {
		if admin == userId {
			return true
		}
	}
	return false
}
