package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var (
	open_ai_api_key        string
	discord_api_key        string
	dev_discord_api_key    string
	elevenlabs_api_key     string
	perplexity_api_key     string
	bot_open_ai_model      string
	open_ai_model          string
	image_generation_model string
	admins                 []string
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
	image_generation_model = os.Getenv("IMAGE_GENERATION_MODEL")
	admins = strings.Split(os.Getenv("ADMINS"), ",")

	// Validate required environment variables
	requiredVars := map[string]string{
		"OPEN_AI_API_KEY":     open_ai_api_key,
		"DISCORD_API_KEY":     discord_api_key,
		"BOT_OPEN_AI_MODEL":   bot_open_ai_model,
		"OPEN_AI_MODEL":       open_ai_model,
		"PERPLEXITY_API_KEY":  perplexity_api_key,
		"IMAGE_GENERATION_MODEL": image_generation_model,
	}

	for name, value := range requiredVars {
		if value == "" {
			log.Fatalf("Required environment variable %s is not set", name)
		}
	}
}

func GetOpenAIKey() string {
	return open_ai_api_key
}

func GetOpenAIModel() string {
	return open_ai_model
}

func GetImageGenerationModel() string {
	return image_generation_model
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
