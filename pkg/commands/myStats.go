package commands

import (
	"fmt"
	"main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func MyStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	userStats, err := persistance.GetUserStats(i.Interaction.Member.User.ID)

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "There was an error getting your stats. Please try again later.",
			},
		})
		return
	}

	userStatsFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Current Streak",
			Value:  fmt.Sprintf("%d day(s)", userStats.BonusStreak),
			Inline: true,
		},
		{
			Name:   "Current Balance",
			Value:  fmt.Sprintf("%.2f token(s)", userStats.ImageTokens),
			Inline: true,
		},
		{}, // Needed for diaplaying next row
		{
			Name:   "Good Bot Count",
			Value:  fmt.Sprintf("%d time(s) praised", userStats.GoodBotCount),
			Inline: true,
		},
		{
			Name:   "Bad Bot Count",
			Value:  fmt.Sprintf("%d time(s) scolded ", userStats.BadBotCount),
			Inline: true,
		},
		{},
		{
			Name:   "Total Images",
			Value:  fmt.Sprintf("%d image(s)", userStats.ImageCount),
			Inline: true,
		},
		{},
		{
			Name:  "Portfolio",
			Value: "To see your portfolio, use /portfolio",
		},
	}

	userStatsEmbed := &discordgo.MessageEmbed{
		Title:       "Your Stats",
		Description: "Here are your current stats",
		Fields:      userStatsFields,
		Color:       0x00FF00,
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{userStatsEmbed},
		},
	})
}
