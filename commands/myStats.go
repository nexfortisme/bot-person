package commands

import (
	"main/persistance"

	"github.com/bwmarrin/discordgo"
)

func MyStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	userStatisticsString := persistance.SlashGetUserStats(*i.Interaction.Member.User)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: userStatisticsString,
		},
	})
}
