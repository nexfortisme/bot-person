package commands

import (
	// "main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type BotStats struct {}

func (b *BotStats) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "bot-stats",
		Description: "Get global usage stats.",
	}
}

func (b *BotStats) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
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

func (b *BotStats) HelpString() string {
	return "The `/bot-stats` command allows you to see global usage stats. This includes the number of interactions with the bot, the number of images requested from the `/image` command, the number of users who have interacted with the bot, and the number of servers the bot is in."
}

func (b *BotStats) CommandCost() int {
	return 0
}
