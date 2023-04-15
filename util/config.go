package util

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

var (
	config ConfigStruct

	createdConfig = false
)

func ReadConfig() {

	var botPersonConfig []byte
	botPersonConfig, err := os.ReadFile("config.json")

	if err != nil {
		createdConfig = true
		log.Printf("Error reading config. Creating File")
		os.WriteFile("config.json", []byte("{\"DiscordToken\":\"\",\"OpenAIKey\":\"\"}"), 0666)
		botPersonConfig, err = os.ReadFile("config.json")
		HandleFatalErrors(err, "Could not read config file: config.json")
	}

	err = json.Unmarshal(botPersonConfig, &config)
	HandleFatalErrors(err, "Could not parse: config.json")

	// Handling the case the config file has just been created
	if config.DiscordToken == "" {
		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter the Discord Token: ")
		config.DiscordToken, _ = reader.ReadString('\n')
		config.DiscordToken = strings.TrimSuffix(config.DiscordToken, "\r\n")
		log.Println("Discord Token Set to: '" + config.DiscordToken + "'")
	}

	// TODO - Check to see if the user doesn't type in a command
	// If they don't, ask them if they wish to continue without OpenAI responses
	if config.OpenAIKey == "" {
		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter the Open AI Key: ")
		config.OpenAIKey, _ = reader.ReadString('\n')
		config.OpenAIKey = strings.TrimSuffix(config.OpenAIKey, "\r\n")
		log.Println("Open AI Key Set to: '" + config.OpenAIKey + "'")
	}

	// TODO - Add check for finnhub token and prompt if the user wants to continue without it

}

func WriteConfig() {
	log.Println("Config Updated. Writing...")

	fle, _ := json.Marshal(config)
	err := os.WriteFile("config.json", fle, 0666)
	if err != nil {
		log.Fatalf("Error Writing config.json")
		return
	}
}

func GetDiscordToken() string {
	return config.DiscordToken
}

func GetDevDiscordToken() string {
	return config.DevDiscordToken
}

func GetOpenAIKey() string {
	return config.OpenAIKey
}

func GetFinHubToken() string {
	return config.FinnHubToken
}

func SetDevDiscordToken(DevDiscordToken string) {
	config.DevDiscordToken = DevDiscordToken
}

func SetFinnHubToken(FinnHubToken string) {
	config.FinnHubToken = FinnHubToken
}
