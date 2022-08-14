package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/util"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

// TODO - Add user tracking
// USERs []USER `json:"users"`
/*
	user: {
		username string
		userData BotTracking
	}
*/
type BotTracking struct {
	BadBotCount  int          `json:"BadBotCount"`
	MessageCount int          `json:"MessageCount"`
	UserStats    []UserStruct `json:"UserTracking"`
}

type UserStruct struct {
	UserName  string          `json:"username"`
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
		ioutil.WriteFile("botTracking.json", []byte("{\"BadBotCount\":0,\"MessageCount\":0}"), 0666)

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

	var hitUser = false
	LogIncomingMessage(s, m, m.Message.Content)

	// TODO - Handle this better. I don't like traversing an array each time.
	for index, element := range botTracking.UserStats {
		fmt.Println("Element Username, Author Username: " + element.UserName + " , " + m.Author.ID)
		if element.UserName != m.Author.ID {
			continue
		} else {
			fmt.Println("User Match Found")
			hitUser = true

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

	if !hitUser {
		fmt.Println("Creating New User For: " + m.Author.Username)
		if flag == 1 {
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{m.Author.ID, UserStatsStruct{0, 1}})
		} else {
			botTracking.UserStats = append(botTracking.UserStats, UserStruct{m.Author.ID, UserStatsStruct{1, 1}})
		}
	}

	if flag == 1 {
		botTracking.MessageCount++
	} else {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	}
}

func GetUserStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, element := range botTracking.UserStats {
		fmt.Println("Element Username, Author Username: " + element.UserName + " , " + m.Author.ID)
		if element.UserName != m.Author.ID {
			continue
		} else {
			msg := "You have interacted with the bot " + strconv.Itoa(element.UserStats.MessageCount) + " times and you scolded the bot " + strconv.Itoa(element.UserStats.BadBotCount) + " times."
			_, err := s.ChannelMessageSend(m.ChannelID, msg)
			util.HandleErrors(err)
			LogOutGoingMessage(s, m, msg)
		}
	}
}

func GetBotStats(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "Across all servers, the bot has been interacted with " + strconv.Itoa(botTracking.MessageCount) + " times and has been bad " + strconv.Itoa(botTracking.BadBotCount) + " times."
	_, err := s.ChannelMessageSend(m.ChannelID, msg)
	util.HandleErrors(err)
	LogOutGoingMessage(s, m, msg)
}

func GetBadBotCount() int {
	return botTracking.BadBotCount
}

func ShutDown() {
	fle, _ := json.Marshal(botTracking)
	ioutil.WriteFile("botTracking.json", fle, 0666)
}
