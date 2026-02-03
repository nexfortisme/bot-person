package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

const (
	defaultBotModel      = "gpt-5-nano"
	defaultImageModel    = "gpt-5"
	defaultOpenAIModel   = "gpt-5-nano"
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

func ReadEnv(useEnvFile bool, devMode bool) {

	var envFilePath string

	if useEnvFile {
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("Tryign to get osEnv: %v", os.Getenv("OPEN_AI_API_KEY"))
			log.Fatalf("Error getting current working directory: %v", err)
		}

		envFilePath = filepath.Join(cwd, ".env")
		
		// Check if .env file exists
		if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
			exampleEnvPath := filepath.Join(cwd, "example.env")
			exampleContent, err := os.ReadFile(exampleEnvPath)
			if err == nil {
				if err := os.WriteFile(envFilePath, exampleContent, 0644); err != nil {
					log.Fatalf("Error creating .env file from example.env: %v", err)
				}
				log.Printf("Created .env file from example.env at %s", envFilePath)
			} else {
				log.Printf(".env file not found at %s and example.env could not be read.", envFilePath)
				log.Println("Create a .env file with the required values, or set environment variables directly.")
			}
		}
		
		err = godotenv.Overload(envFilePath)
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
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

	if bot_open_ai_model == "" {
		bot_open_ai_model = defaultBotModel
	}
	if image_generation_model == "" {
		image_generation_model = defaultImageModel
	}
	if open_ai_model == "" {
		open_ai_model = defaultOpenAIModel
	}

	// Validate required environment variables
	requiredVars := map[string]string{
		"OPEN_AI_API_KEY":    open_ai_api_key,
		"DISCORD_API_KEY":    discord_api_key,
		"OPEN_AI_MODEL":      open_ai_model,
		"PERPLEXITY_API_KEY": perplexity_api_key,
	}

	if devMode {
		requiredVars["DEV_DISCORD_API_KEY"] = dev_discord_api_key
	}

	var missing []string
	for name, value := range requiredVars {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		log.Printf("Missing required configuration: %s", strings.Join(missing, ", "))
		if useEnvFile {
			if envFilePath == "" {
				log.Println("Update your .env file or set environment variables directly, then retry.")
			} else {
				log.Printf("Update %s or set environment variables directly, then retry.", envFilePath)
			}
		} else {
			log.Println("Set environment variables directly, or run with --useEnvFile to load from .env.")
		}
		log.Fatal("Configuration incomplete.")
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
