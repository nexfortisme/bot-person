package commands

import (
	"fmt"
	"main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func MyStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	userStats, _ := persistance.GetUserStats(i.Interaction.Member.User.ID)

	myStatsEmbed := &discordgo.MessageEmbed{
		Title:       "Your Stats",
		Description: "Here are your stats!",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			// {
			// 	Name:  "Total Interactions",
			// 	Value: fmt.Sprintf("%d Interaction(s)", userStats.),
			// },
			// {},
			{
				Name:  "Times Praised Bot Person",
				Value: fmt.Sprintf("%d Time(s)", userStats.GoodBotCount),
				Inline: true,
			},
			{
				Name:  "Times Scolded Bot Person",
				Value: fmt.Sprintf("%d Time(s)", userStats.BadBotCount),
				Inline: true,
			},
			{},
			{
				Name:  "Image(s) Requested",
				Value: fmt.Sprintf("%d Image(s)", userStats.ImageCount),
				Inline: false,
			},
			{},
			{
				Name:  "Bonus Streak",
				Value: fmt.Sprintf("%d Day(s)", userStats.BonusStreak),
				Inline: true,
			},
			{
				Name:  "Token Balance",
				Value: fmt.Sprintf("%.2f Token(s)", userStats.ImageTokens),
				Inline: true,
			},
		},
	}

	// Getting user stat data
	// userStatisticsString := persistance.SlashGetUserStats(*i.Interaction.Member.User)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{myStatsEmbed},
			// Components: []discordgo.MessageComponent{actionRow},
		},
	})
}
