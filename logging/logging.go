package logging

import (
	"log"
	"main/util"

	"github.com/bwmarrin/discordgo"
)

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
