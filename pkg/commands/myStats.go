package commands

import (
	// persistance "main/pkg/persistance/services"

	// logging "main/pkg/logging/services"
	// loggingType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func MyStats(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// userStats, _ := persistance.GetUserStats(i.Interaction.Member.User.ID, s)

	// userStatsString :=

	// persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	// userStatisticsString := persistance.SlashGetUserStats(*i.Interaction.Member.User)

	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: userStatisticsString,
	// 	},
	// })
}
