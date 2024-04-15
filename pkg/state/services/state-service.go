package state

import (
	"github.com/bwmarrin/discordgo"
)

var (
	discordSession *discordgo.Session
)

func SetDiscordSession(s *discordgo.Session) {
	discordSession = s
}

func GetDiscordSession() *discordgo.Session {
	return discordSession
}
