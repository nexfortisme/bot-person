package messages

import (
	"log"

	"github.com/bwmarrin/discordgo"
)


func ParseMessage() {

}

func InitLogging(){

}

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)
}