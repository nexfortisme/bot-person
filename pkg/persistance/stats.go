package persistance

import (
	"fmt"
	"log"
	persistance "main/pkg/persistance/models"

	"github.com/bwmarrin/discordgo"
)

type BPInteraction int

const (
	BPChatInteraction       BPInteraction = iota // anything with the divinci chat model
	BPImageRequest                               // any image request
	BPBadBotInteraction                          // bad bot
	BPGoodBotInteraction                         // good bot
	BPBasicInteraction                           // any basic bot interaction
	BPSystemInteraction                          // TODO - Figure out what this means
	BPCreateUserInteraction                      // used by system for creating users
	BPLennyFaceInteracton
)

// func IncrementInteractionTracking(flag BPInteraction, user discordgo.User) {

// 	userId := user.ID
// 	username := user.Username

// 	foundUser := false
// 	foundUser = handleUserStatIncrementing(flag, userId)

// 	if !foundUser {
// 		createNewUserTracking(flag, userId, username)
// 	}
// }

func printUserStocks(user persistance.User) string {

	if len(user.UserStats.Stocks) == 0 {
		return "You don't have any stocks."
	}

	retString := "You have the following stocks:\n"

	for _, element := range user.UserStats.Stocks {
		retString += fmt.Sprintf("\t%s: %.2f\n", element.StockTicker, element.StockCount)
	}

	return retString
}

// func SlashGetBotStats(s *discordgo.Session) string {

// 	guildCount := len(s.State.Guilds)

// 	var globalMessageCount int = 0
// 	var globalGoodBotCount int = 0
// 	var globalBadBotCount int = 0
// 	var globalImageCount int = 0
// 	// var globalLongestBonusStreak int = 0

// 	var globalTokenCirculation float64 = 0.0

// 	var returnMessage string

// 	for _, element := range botTracking.UserStats {

// 		globalMessageCount += element.UserStats.MessageCount
// 		globalGoodBotCount += element.UserStats.GoodBotCount
// 		globalBadBotCount += element.UserStats.BadBotCount
// 		globalImageCount += element.UserStats.ImageCount

// 		globalTokenCirculation += element.UserStats.ImageTokens
// 	}

// 	returnMessage = fmt.Sprintf("Across %d servers, Bot Person has/is/did:\nInteractions: %d\nBeen Good: %d\nBeen Bad: %d\nGenerated Images: %d\nTotal Tokens In Circulation: %.2f", guildCount, globalMessageCount, globalGoodBotCount, globalGoodBotCount, globalImageCount, globalTokenCirculation)

// 	return returnMessage
// }

func PrintUSerStocksHelper(user discordgo.User) (string, error) {
	userStruct, err := GetUser(user.ID)

	if err != nil {
		log.Println("Error getting user: " + err.Error())
		return "", err
	}

	return printUserStocks(*userStruct), nil
}
