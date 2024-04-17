package commands

import (
	// "main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func BotStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.COMMAND_BOT_STATS, i.Interaction.Member.User.ID, "Bot Stats command used", i.Interaction.GuildID)

	// Getting user stat data
	botStatisticsString := "Refactor in progress..."

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: botStatisticsString,
		},
	})
}