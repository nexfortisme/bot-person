package commands

import (
	"main/persistance"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	helpString := "A picture is worth 1000 words"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpString,
		},
	})
}