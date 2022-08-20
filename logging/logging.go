package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/util"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type BotTracking struct {
	BadBotCount  int          `json:"BadBotCount"`
	MessageCount int          `json:"MessageCount"`
	UserStats    []UserStruct `json:"UserTracking"`
}

type UserStruct struct {
	UserId    string          `json:"username"`
	UserStats UserStatsStruct `json:"userStats"`
}

type UserStatsStruct struct {
	BadBotCount  int `json:"BadBotCount"`
	MessageCount int `json:"MessageCount"`
}

var (
	botTracking BotTracking
)

func InitBotStatistics() {
	var trackingFile []byte

	trackingFile, err := ioutil.ReadFile("botTracking.json")
	if err != nil {

		log.Printf("Error Reading botTracking. Creating File")
		os.WriteFile("botTracking.json", []byte("{\"BadBotCount\":0,\"MessageCount\":0, \"UserStats\":[]}"), 0666)

		trackingFile, err = ioutil.ReadFile("botTracking.json")
		if err != nil {
			log.Fatalf("Could not read config file: botTracking.json")
		}

	}

	err = json.Unmarshal(trackingFile, &botTracking)
	if err != nil {
		log.Fatalf("Could not parse: botTracking.json")
	}
}

func LogOutGoingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {

}

func LogError(err string) {
	log.Fatalf(err)
}

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)

}

func IncrementTracker(flag int, m *discordgo.MessageCreate, s *discordgo.Session) {

	var foundUser = false
	LogIncomingMessage(s, m, util.ReplaceIDsWithNames(m, s))

	// TODO - Handle this better. I don't like traversing an array each time.
	for index, element := range botTracking.UserStats {
		if element.UserId != m.Author.ID {
			continue
		} else {
			foundUser = true

			if flag == 1 {
				element.UserStats.MessageCount++
			} else {
				element.UserStats.MessageCount++
				element.UserStats.BadBotCount++
			}

			// Is this necessary?
			botTracking.UserStats[index] = element
		}
	}

	if !foundUser {
		createNewUserTracking(m.Author.ID, m.Author.Username, flag)
	}

	incrementBotTracking(flag)
}

func IncreametSlashCommandTracker(flag int, userId string, username string) {
	foundUser := false

	for index, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			foundUser = true

			if flag == 1 {
				element.UserStats.MessageCount++
			} else {
				element.UserStats.MessageCount++
				element.UserStats.BadBotCount++
			}

			botTracking.UserStats[index] = element
		}
	}

	if !foundUser {
		createNewUserTracking(userId, username, flag)
	}

	incrementBotTracking(flag)
}

func GetUserStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, element := range botTracking.UserStats {
		fmt.Println("Element Username, Author Username: " + element.UserId + " , " + m.Author.ID)
		if element.UserId != m.Author.ID {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times and you scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times."
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
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times and you scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times."
			return msg
		}
	}
	msg := "Sorry, you don't have any recorded interactions with the bot."

	return msg
}

func GetBotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func SlashGetBotStats(s *discordgo.Session) string{
	guildCount := len(s.State.Guilds)
	msg := "Across " + strconv.Itoa(guildCount) + " servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	return msg
}

func GetBadBotCount() int {
	return botTracking.BadBotCount
}

func ShutDown() {
	fle, _ := json.Marshal(botTracking)
	os.WriteFile("botTracking.json", fle, 0666)
}

func incrementBotTracking(flag int) {
	if flag == 1 {
		botTracking.MessageCount++
	} else {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	}
}

func createNewUserTracking(userId string, username string, flag int){
	log.Println("Creating New User For: " + username)
	if flag == 1 {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{0, 1}})
	} else {
		botTracking.UserStats = append(botTracking.UserStats, UserStruct{userId, UserStatsStruct{1, 1}})
	}
}
