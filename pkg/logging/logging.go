package logging

import (
	"log"
	"main/pkg/util"

	"github.com/bwmarrin/discordgo"
)

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name
	message := util.ReplaceIDsWithNames(m, s)

	log.Printf("%s (%s) > %s\n", requestUser, rGuildName, message)
}

func LogIncomingUserInteraction(s *discordgo.Session, requestUser string, guildId string, message string) {
	rGuild, _ := s.State.Guild(guildId)
	rGuildName := rGuild.Name

	log.Printf("%s (%s) > %s\n", requestUser, rGuildName, message)
}

func LogOutgoingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name
	message := util.ReplaceIDsWithNames(m, s)

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)
}

func LogOutgoingUserInteraction(s *discordgo.Session, requestUser string, guildId string, message string) {
	rGuild, _ := s.State.Guild(guildId)
	rGuildName := rGuild.Name

	log.Printf("%s (%s) < %s\n", requestUser, rGuildName, message)
}

func LogError(err string) {
	log.Fatalf(err)
}
