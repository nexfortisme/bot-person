package logging

import (
	"log"
	"main/util"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func IncrementTracker(flag int, userId string, username string) {
	foundUser := false
	foundUser = handleUserStatIncrementing(flag, userId)

	if !foundUser {
		createNewUserTracking(userId, username, flag)
	}

	incrementBotTracking(flag)
}

func GetUserStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, element := range botTracking.UserStats {
		if element.UserId != m.Author.ID {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times, praised the bot " + strconv.Itoa(element.UserStats.GoodBotCount) + " times, and scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times. You have requested an image " + strconv.Itoa(element.UserStats.ImageCount) + " times."
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			util.HandleErrors(err)
			LogOutGoingMessage(s, m, msg)
			return
		}
	}

	msg := "Sorry, you don't have any recorded interactions with the bot."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func SlashGetUserStats(s *discordgo.Session, i *discordgo.InteractionCreate) string {
	userId := i.Interaction.Member.User.ID

	for _, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times, praised the bot " + strconv.Itoa(element.UserStats.GoodBotCount) + " times, and scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times. You have requested an image " + strconv.Itoa(element.UserStats.ImageCount) + " times."
			return msg
		}
	}
	msg := "Sorry, you don't have any recorded interactions with the bot."

	return msg
}

func GetBadBotCount() int {
	return botTracking.BadBotCount
}

func GetBotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times, praised " + strconv.Itoa(botTracking.GoodBotCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func SlashGetBotStats(s *discordgo.Session) string {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times, praised " + strconv.Itoa(botTracking.GoodBotCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	return msg
}

func incrementBotTracking(flag int) {
	if flag == 1 {
		botTracking.MessageCount++
		botTracking.GoodBotCount++
	} else if flag == 2 {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	} else {
		botTracking.MessageCount++
	}
}

func createNewUserTracking(userId string, username string, flag int) {
	log.Println("Creating New User For: " + username)
	if flag == 1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 1, 0, 0, 0}})
	} else if flag == 2 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 1, 0, 0}})
	} else if flag == 3 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 1, 1, 0}})
	} else {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 0, 0, 0, 0}})
	}
}

func handleUserStatIncrementing(flag int, userId string) bool {

	// TODO - Handle this better. I don't like traversing an array each time.
	// Convert to a map....eventually
	for index, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {

			if flag == 1 {
				element.UserStats.MessageCount++
				element.UserStats.GoodBotCount++
			} else if flag == 2 {
				element.UserStats.MessageCount++
				element.UserStats.BadBotCount++
			} else if flag == 3 {
				element.UserStats.MessageCount++
				element.UserStats.ImageCount++
			} else {
				element.UserStats.MessageCount++
			}

			// Is this necessary?
			botTracking.UserStats[index] = element
			return true
		}
	}

	return false
}
