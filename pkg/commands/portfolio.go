package commands

import (
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Portfolio(s *discordgo.Session, i *discordgo.InteractionCreate) {

	logging.LogEvent(eventType.COMMAND_PORTFOLIO, i.Interaction.Member.User.ID, "Portfolio command used", i.Interaction.GuildID)

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
