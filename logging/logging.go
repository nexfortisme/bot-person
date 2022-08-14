package logging

import (
	"encoding/json"
	"io/ioutil"
	"log"

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

func LogError(err string){
	log.Fatalf(err);
}

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)

}

func IncrementTracker(flag int) {
	if flag == 1 {
		botTracking.MessageCount++
	} else {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	}
}



func GetBadBotCount() int {
	return botTracking.BadBotCount
}

func ShutDown() {
	fle, _ := json.Marshal(botTracking)
	ioutil.WriteFile("botTracking.json", fle, 0666)
}