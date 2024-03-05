package commands

import (
	// "main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func BotStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	// botStatisticsString := persistance.SlashGetBotStats(s)

	// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 	Data: &discordgo.InteractionResponseData{
	// 		Content: botStatisticsString,
	// 	},
	// })
}