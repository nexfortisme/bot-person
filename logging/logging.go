package logging

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"main/util"
	"os"

	"github.com/bwmarrin/discordgo"
)

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

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name
	message := util.ReplaceIDsWithNames(m, s)

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)
}

func ShutDown() {
	log.Println("Writing botTracking.json...")
	fle, _ := json.Marshal(botTracking)
	os.WriteFile("botTracking.json", fle, 0666)
}

func getUser(userId string) (UserStatsStruct, error) {
	for _, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {
			return element.UserStats, nil
		}
	}
	return UserStatsStruct{-1, -1, -1, -1, -1}, nil
}

func updateUser(userStats UserStatsStruct, userId string) bool {
	for index, element := range botTracking.UserStats {
		if element.UserId != userId {
			continue
		} else {

			// Is this necessary?
			botTracking.UserStats[index].UserStats = userStats
			return true
		}
	}
	return false
}
