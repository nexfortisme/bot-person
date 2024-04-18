package commands

import (
	"fmt"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func MyStats(s *discordgo.Session, i *discordgo.InteractionCreate) {

	logging.LogEvent(eventType.COMMAND_MY_STATS, i.Interaction.Member.User.ID, "My Stats command used", i.Interaction.GuildID)

	user, _ := persistance.GetUser(i.Interaction.Member.User.ID)

	myStatsEmbed := &discordgo.MessageEmbed{
		Title:       "Your Stats",
		Description: user.ID,
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			// {
			// 	Name:  "Total Interactions",
			// 	Value: fmt.Sprintf("%d Interaction(s)", userStats.),
			// },
			// {},
			// {
			// 	Name:   "Times Praised Bot Person",
			// 	Value:  fmt.Sprintf("%d Time(s)", user.UserStats.GoodBotCount),
			// 	Inline: true,
			// },
			// {
			// 	Name:   "Times Scolded Bot Person",
			// 	Value:  fmt.Sprintf("%d Time(s)", user.UserStats.BadBotCount),
			// 	Inline: true,
			// },
			{},
			{
				Name:   "Image(s) Requested",
				Value:  fmt.Sprintf("%d Image(s)", user.UserStats.ImageCount),
				Inline: false,
			},
			{},
			{
				Name:   "Bonus Streak",
				Value:  fmt.Sprintf("%d Day(s)", user.UserStats.BonusStreak),
				Inline: true,
			},
			{
				Name:   "Token Balance",
				Value:  fmt.Sprintf("%.2f Token(s)", user.UserStats.ImageTokens),
				Inline: true,
			},
		},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{myStatsEmbed},
			// Components: []discordgo.MessageComponent{actionRow},
		},
	})
}
