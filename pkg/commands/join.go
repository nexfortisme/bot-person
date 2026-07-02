package commands

import (
	// "main/pkg/external"
	// state "main/pkg/state/services"

	"github.com/bwmarrin/discordgo"
)

func Join(s *discordgo.Session, i *discordgo.InteractionCreate) {

	channel, _ := s.Channel(i.Interaction.ChannelID)

	// Check to see if channel is a voice channel
	if channel.Type == discordgo.ChannelTypeGuildVoice {

		// Join the user's voice channel
		
		// state.SetConnection(i.Message.ChannelID, external.Join(s, i.Interaction))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Bot Joined Channel",
			},
		})
		return
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You must be in a voice channel to use this command",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

}
