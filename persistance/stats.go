package persistance

import (
	"log"
	"strconv"

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
		return "You have interacted with the bot " + strconv.Itoa(bpUser.UserStats.MessageCount) + " times, praised the bot " + strconv.Itoa(bpUser.UserStats.GoodBotCount) + " times, and scolded the bot " + strconv.Itoa(bpUser.UserStats.BadBotCount) + " times. You have requested an image " + strconv.Itoa(bpUser.UserStats.ImageCount) + " times."
	}
}

// TODO - Rewrite using new global stat tracking
func SlashGetBotStats(s *discordgo.Session) string {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times, praised " + strconv.Itoa(botTracking.GoodBotCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	return msg
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
