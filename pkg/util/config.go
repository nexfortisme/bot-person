package util

import (
	"bufio"
	"encoding/json"
	"log"
	util "main/pkg/util/models"
	"os"
	"strings"
)

var (
	config util.Config

	createdConfig = false
)

func ReadConfig() {

	var botPersonConfig []byte
	botPersonConfig, err := os.ReadFile("config.json")

	// Going to assume that if there is an error, it is because the file doesn't exist
	if err != nil {
		createdConfig = true

		log.Printf("Error reading config. Creating File")
		_, err = os.Create("config.json")

		botPersonConfig, err = os.ReadFile("config.json")
	}

	err = json.Unmarshal(botPersonConfig, &config)

	if config.DiscordToken == "" {
		createdConfig = true
		readAPIKey(&config.DiscordToken, "Discord Token")
	}

	if config.OpenAIKey == "" {
		createdConfig = true
		readAPIKey(&config.OpenAIKey, "Open AI Key")
	}

	if config.FinnHubToken == "" {
		createdConfig = true
		readAPIKey(&config.FinnHubToken, "FinnHub Token")
	}

	if config.AdminIDs == nil {
		createdConfig = true

		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter an Admin ID: ")
		adminID, _ := reader.ReadString('\n')
		adminID = strings.TrimSuffix(adminID, "\r\n")

		config.AdminIDs = append(config.AdminIDs, adminID)

		log.Printf("Admin ID Added: '%s'\n", adminID)
	}

	if createdConfig {
		WriteConfig()
	}

}

func readAPIKey(variable *string, flavorText string) {
	reader := bufio.NewReader(os.Stdin)
	log.Printf("Please Enter the %s: ", flavorText)
	*variable, _ = reader.ReadString('\n')
	*variable = strings.TrimSuffix(*variable, "\r\n")
	log.Printf("%s Set to: '%s'\n", flavorText, *variable)
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

func GetAdminIds() []string {
	return config.AdminIDs
}

func AddAdmin(userId string) {
	config.AdminIDs = append(config.AdminIDs, userId)
}

func RemoveAdmin(userId string) {
	for i, id := range config.AdminIDs {
		if id == userId {
			config.AdminIDs = append(config.AdminIDs[:i], config.AdminIDs[i+1:]...)
		}
	}
}

func UserIsAdmin(userId string) bool {
	for _, id := range config.AdminIDs {
		if strings.Compare(id, userId) == 0 {
			return true
		}
	}
	return false
}

func ListAdmins() string {
	var adminList string = ""
	for index, id := range config.AdminIDs {
		if index == len(config.AdminIDs)-1 {
			adminList += id
			break
		}
		adminList += id + ", "
	}
	return adminList
}
