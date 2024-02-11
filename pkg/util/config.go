package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	db_host     string
	db_user     string
	db_password string
	db_name     string

	open_ai_api_key     string
	discord_api_key     string
	dev_discord_api_key string
	finn_hub_api_key    string
)

func ReadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db_host = os.Getenv("DB_HOST")
	db_user = os.Getenv("DB_USER")
	db_password = os.Getenv("DB_PASSWORD")
	db_name = os.Getenv("DB_NAME")

	open_ai_api_key = os.Getenv("OPENAI_API_KEY")
	discord_api_key = os.Getenv("DISCORD_API_KEY")
	dev_discord_api_key = os.Getenv("DEV_DISCORD_API_KEY")
	finn_hub_api_key = os.Getenv("FINNHUB_API_KEY")

	fmt.Println("DB_HOST: ", db_host)
	fmt.Println("DB_USER: ", db_user)
	fmt.Println("DB_PASSWORD: ", db_password)
	fmt.Println("DB_NAME: ", db_name)

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

func GetDBName() string {
	return db_name
}

func GetOpenAIKey() string {
	return open_ai_api_key
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

func UserIsAdmin(userId string) bool {
	return false
}

// func AddAdmin(userId string) {
// 	config.AdminIDs = append(config.AdminIDs, userId)
// }

// func RemoveAdmin(userId string) {
// 	for i, id := range config.AdminIDs {
// 		if id == userId {
// 			config.AdminIDs = append(config.AdminIDs[:i], config.AdminIDs[i+1:]...)
// 		}
// 	}
// }

// func UserIsAdmin(userId string) bool {
// 	for _, id := range config.AdminIDs {
// 		if strings.Compare(id, userId) == 0 {
// 			return true
// 		}
// 	}
// 	return false
// }

// func ListAdmins() string {
// 	var adminList string = ""
// 	for index, id := range config.AdminIDs {
// 		if index == len(config.AdminIDs)-1 {
// 			adminList += id
// 			break
// 		}
// 		adminList += id + ", "
// 	}
// 	return adminList
// }
