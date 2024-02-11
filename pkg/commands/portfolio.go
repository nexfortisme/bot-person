package commands

import (
	"main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func Portfolio(s *discordgo.Session, i *discordgo.InteractionCreate) {

	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	userStatisticsString, err := persistance.PrintUSerStocksHelper(*i.Interaction.Member.User)

	if err != nil {
		userStatisticsString = "Something went wrong"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: userStatisticsString,
		},
	})

}
