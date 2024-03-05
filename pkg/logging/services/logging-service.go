package logging

import (
	"log"
	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/models"
	"main/pkg/persistance"
	"main/pkg/util"

	"github.com/bwmarrin/discordgo"
)

func LogEvent(event loggingType.EventType, eventValue string, createUser string, createGuildId string, session *discordgo.Session) (bool, error) {

	db := persistance.GetDB()

	var guild *discordgo.Guild

	if session == nil {
		guild.Name = "System"
	} else {
		guild, _ = session.State.Guild(createGuildId)
	}

	loggingEvent := logging.LoggingEvent{
		EventType:     event,
		EventValue:    eventValue,
		CreateGuildId: createGuildId,
		CreateGuild:   guild.Name,
	}

	_, err := db.Model(&loggingEvent).Insert()
	if err != nil {
		return false, err
	}

	return true, nil
}

func LogIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name
	message := util.ReplaceIDsWithNames(m, s)

	log.Printf("%s (%s) > %s\n", requestUser, rGuildName, message)
}

func GetEventUserIdAndGuild(s *discordgo.Session, m *discordgo.MessageCreate) (string, string) {
	requestUser := m.Author.ID
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name

	return requestUser, rGuildName
}
