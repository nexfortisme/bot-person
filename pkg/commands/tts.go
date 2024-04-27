package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func TTS(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Find the guild for the message
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Error finding guild")
		return
	}

	// Find the voice state for the user
	for _, vs := range guild.VoiceStates {
		if vs.UserID == m.Author.ID {
			// Join the user's voice channel
			vc, err := s.ChannelVoiceJoin(guild.ID, vs.ChannelID, false, true)
			if err != nil {
				fmt.Println("Error joining the voice channel:", err)
				return
			}
			fmt.Println("Joined voice channel:", vc.ChannelID)
			break
		}
	}
}

func Leave(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Find the guild for the message
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		fmt.Println("Error finding guild")
		return
	}

	fmt.Printf("Voice Channels: %+v\n", s.VoiceConnections)

	connections := s.VoiceConnections[guild.ID]

	// Disconnect from the voice channel
	err = connections.Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting from the voice channel:", err)
	} else {
		fmt.Println("Disconnected from the voice channel")
	}

}
