package commands

import (
	"main/persistance"

	"github.com/bwmarrin/discordgo"
)

func Broken(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	brokenMessage := "If you have something that is broken about Bot Person, you can create an issue describing what you found here: https://github.com/nexfortisme/bot-person/issues/new"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: brokenMessage,
		},
	})
}
