package state

import (
	"github.com/bwmarrin/discordgo"
)

var (
	discordSession *discordgo.Session
	connections = make(map[string]*discordgo.VoiceConnection)
)

func SetDiscordSession(s *discordgo.Session) {
	discordSession = s
}

func GetDiscordSession() *discordgo.Session {
	return discordSession
}

func SetConnection(channelId string, vc *discordgo.VoiceConnection) {
	connections[channelId] = vc
}

func GetConnection(channelId string) *discordgo.VoiceConnection {
	return connections[channelId]
}

func RemoveConnection(channelId string) {
	connections[channelId] = nil
}