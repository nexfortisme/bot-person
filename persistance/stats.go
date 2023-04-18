package persistance

import (
	"fmt"
	"log"

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

func IncrementInteractionTracking(flag BPInteraction, user discordgo.User) {

	userId := user.ID
	username := user.Username

	foundUser := false
	foundUser = handleUserStatIncrementing(flag, userId)

	if !foundUser {
		createNewUserTracking(flag, userId, username)
	}
}

func SlashGetUserStats(user discordgo.User) string {
	bpUser, err := getUser(user.ID)

	if err != nil {
		return "Sorry, you don't have any recorded interactions with the bot."
	} else {
		return fmt.Sprintf("You have interacted with the bot %d times.\nYou praised the bot %d times and scolded the bot %d times.\nYou have requested %d images.\nYour current bonus streak is %d.\nYou have %d Save Streak Tokens.\n%s", bpUser.UserStats.MessageCount, bpUser.UserStats.GoodBotCount, bpUser.UserStats.BadBotCount, bpUser.UserStats.ImageCount, bpUser.UserStats.BonusStreak, bpUser.UserStats.SaveStreakTokens, printUserStocks(bpUser))
	}
}

func printUserStocks(user UserStruct) string {

	if len(user.UserStats.Stocks) == 0 {
		return "You don't have any stocks."
	}

	retString := "You have the following stocks:\n"

	for _, element := range user.UserStats.Stocks {
		retString += fmt.Sprintf("\t%s: %.2f\n", element.StockTicker, element.StockCount)
	}

	return retString
}

func createNewUserTracking(flag BPInteraction, userId string, username string) {
	log.Println("Creating New User For: " + username)

	switch interaction := flag; interaction {
	case BPChatInteraction:
	case BPLennyFaceInteracton:
	case BPBasicInteraction:
		createAndAddUser(userId, 1, 0, 0, 0, 25)
	case BPImageRequest:
		createAndAddUser(userId, 1, 0, 0, 1, 25)
	case BPBadBotInteraction:
		createAndAddUser(userId, 1, 0, 1, 0, 25)
	case BPGoodBotInteraction:
		createAndAddUser(userId, 1, 1, 0, 0, 25)
	default:
		createAndAddUser(userId, 1, 0, 0, 0, 25)
	}
}

func handleUserStatIncrementing(flag BPInteraction, userId string) bool {

	incrementUser, err := getUser(userId)

	if err != nil {
		return false
	}

	switch interaction := flag; interaction {
	case BPChatInteraction:
	case BPLennyFaceInteracton:
	case BPBasicInteraction:
		incrementUser.UserStats.MessageCount++
	case BPImageRequest:
		incrementUser.UserStats.ImageCount++
	case BPBadBotInteraction:
		incrementUser.UserStats.MessageCount++
		incrementUser.UserStats.BadBotCount++
	case BPGoodBotInteraction:
		incrementUser.UserStats.MessageCount++
		incrementUser.UserStats.GoodBotCount++
	default:
		incrementUser.UserStats.MessageCount++
	}

	updateUser(incrementUser)
	return true
}
