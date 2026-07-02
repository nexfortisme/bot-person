package util

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	open_ai_api_key        string
	discord_api_key        string
	dev_discord_api_key    string
	elevenlabs_api_key     string
	perplexity_api_key     string
	bot_open_ai_model                 string
	open_ai_model                     string
	image_generation_model            string
	local_llm_model                   string
	retort_machine_model              string
	local_llm_chat_completions_endpoint string
	admins                            []string

	modelMu sync.RWMutex

	// Set at startup from process environment before .env is loaded.
	// ReloadModelConfig restores these so Docker/compose overrides are not
	// overwritten by localhost values in the mounted .env file.
	preservedFromRuntimeEnv map[string]string
)

const (
	defaultLocalLLMModel                  = "gemma4:12b"
	defaultRetortMachineModel             = "gemma4:12b"
	defaultLocalLLMChatCompletionsEndpoint = "http://localhost:11434/v1/chat/completions"
)

func ReadEnv(useEnvFile bool) {

	if useEnvFile {
		cwd, err := os.Getwd()
		if err != nil {

			log.Printf("Tryign to get osEnv: %v", os.Getenv("OPEN_AI_API_KEY"))

			log.Fatalf("Error getting current working directory: %v", err)
		}

		preservedFromRuntimeEnv = capturePreservedEnv()
		envFilePath := filepath.Join(cwd, ".env")
		err = loadEnvFile(envFilePath)
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
	local_llm_model = envOrDefault("LOCAL_LLM_MODEL", defaultLocalLLMModel)
	retort_machine_model = envOrDefault("RETORT_MACHINE_MODEL", defaultRetortMachineModel)
	local_llm_chat_completions_endpoint = envOrDefault("LOCAL_LLM_CHAT_COMPLETIONS_ENDPOINT", defaultLocalLLMChatCompletionsEndpoint)
	admins = strings.Split(os.Getenv("ADMINS"), ",")

	// Validate required environment variables
	requiredVars := map[string]string{
		"OPEN_AI_API_KEY":        open_ai_api_key,
		"DISCORD_API_KEY":        discord_api_key,
		"BOT_OPEN_AI_MODEL":      bot_open_ai_model,
		"OPEN_AI_MODEL":          open_ai_model,
		"PERPLEXITY_API_KEY":     perplexity_api_key,
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
	modelMu.RLock()
	defer modelMu.RUnlock()
	return open_ai_model
}

func GetImageGenerationModel() string {
	modelMu.RLock()
	defer modelMu.RUnlock()
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
	modelMu.RLock()
	defer modelMu.RUnlock()
	return bot_open_ai_model
}

func GetLocalLLMModel() string {
	modelMu.RLock()
	defer modelMu.RUnlock()
	return local_llm_model
}

func GetRetortMachineModel() string {
	modelMu.RLock()
	defer modelMu.RUnlock()
	return retort_machine_model
}

func GetLocalLLMChatCompletionsEndpoint() string {
	modelMu.RLock()
	defer modelMu.RUnlock()
	return local_llm_chat_completions_endpoint
}

// ReloadModelConfig re-reads model variables from the .env file (or environment).
// Safe to call from any goroutine. Does not affect API keys or Discord tokens.
func ReloadModelConfig(useEnvFile bool) {
	if useEnvFile {
		cwd, err := os.Getwd()
		if err != nil {
			log.Printf("ReloadModelConfig: error getting working directory: %v", err)
			return
		}
		envFilePath := filepath.Join(cwd, ".env")
		if err := loadEnvFile(envFilePath); err != nil {
			log.Printf("ReloadModelConfig: error loading .env file: %v", err)
			return
		}
	}

	newBot := os.Getenv("BOT_OPEN_AI_MODEL")
	newModel := os.Getenv("OPEN_AI_MODEL")
	newImage := os.Getenv("IMAGE_GENERATION_MODEL")
	newLocalLLM := envOrDefault("LOCAL_LLM_MODEL", defaultLocalLLMModel)
	newRetortMachine := envOrDefault("RETORT_MACHINE_MODEL", defaultRetortMachineModel)
	newLocalLLMEndpoint := envOrDefault("LOCAL_LLM_CHAT_COMPLETIONS_ENDPOINT", defaultLocalLLMChatCompletionsEndpoint)

	modelMu.Lock()
	bot_open_ai_model = newBot
	open_ai_model = newModel
	image_generation_model = newImage
	local_llm_model = newLocalLLM
	retort_machine_model = newRetortMachine
	local_llm_chat_completions_endpoint = newLocalLLMEndpoint
	modelMu.Unlock()

	log.Printf(
		"ReloadModelConfig: models updated — bot=%s, chat=%s, image=%s, local_llm=%s, retort_machine=%s, local_llm_endpoint=%s",
		newBot, newModel, newImage, newLocalLLM, newRetortMachine, newLocalLLMEndpoint,
	)
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Keys listed here may be set by Docker Compose (or the shell) before startup.
// Values present in the process environment are preserved across .env reloads.
var preservedEnvKeys = []string{
	"LOCAL_LLM_CHAT_COMPLETIONS_ENDPOINT",
}

func capturePreservedEnv() map[string]string {
	preserved := make(map[string]string, len(preservedEnvKeys))
	for _, key := range preservedEnvKeys {
		if value := os.Getenv(key); value != "" {
			preserved[key] = value
		}
	}
	return preserved
}

func loadEnvFile(path string) error {
	if err := godotenv.Overload(path); err != nil {
		return err
	}
	restorePreservedEnv(preservedFromRuntimeEnv)
	return nil
}

func restorePreservedEnv(preserved map[string]string) {
	for key, value := range preserved {
		os.Setenv(key, value)
	}
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
