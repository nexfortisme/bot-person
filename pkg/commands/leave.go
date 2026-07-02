package commands

import (
	"main/pkg/external"
	state "main/pkg/state/services"

	"github.com/bwmarrin/discordgo"
)

func Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	
	external.Leave(state.GetConnection(i.Interaction.ChannelID));
	state.RemoveConnection(i.Interaction.ChannelID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Bot Left Channel",
		},
	})
}
