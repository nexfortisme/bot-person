package external

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Join(s *discordgo.Session, i *discordgo.Interaction) *discordgo.VoiceConnection {

	var vc *discordgo.VoiceConnection
	var channel, err = s.Channel(i.ChannelID)

	// Check to see if channel is a voice channel
	if channel.Type != discordgo.ChannelTypeGuildVoice {
		fmt.Println("Channel is not a voice channel")
		return nil
	} else {
		vc, err = s.ChannelVoiceJoin(i.GuildID, i.ChannelID, false, true)
		if err != nil {
			fmt.Println("Error joining the voice channel:", err)
			return nil
		}
		fmt.Println("Joined voice channel:", vc.ChannelID)
	}

	return vc
}

func Leave(vc *discordgo.VoiceConnection) {

	if vc == nil {
		return
	}

	// Disconnect from the voice channel
	err := vc.Disconnect()
	if err != nil {
		fmt.Println("Error disconnecting from the voice channel:", err)
	} else {
		fmt.Println("Disconnected from the voice channel")
	}

}
