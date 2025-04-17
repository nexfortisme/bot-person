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
	userStats := persistance.GetUserStats(i.Interaction.Member.User.ID)

	myStatsEmbed := &discordgo.MessageEmbed{
		Title:       "Your Stats",
		Description: fmt.Sprintf("User ID: %s", user.ID),
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Total Interactions",
				Value:  fmt.Sprintf("%d Interaction(s)", userStats.InteractionCount),
				Inline: true,
			},
			{
				Name:   "Chats with Bot Person",
				Value:  fmt.Sprintf("%d Chat(s)", userStats.ChatCount),
				Inline: true,
			},
			{},
			{
				Name:   "Times Praised Bot Person",
				Value:  fmt.Sprintf("%d Time(s)", userStats.GoodBotCount),
				Inline: true,
			},
			{
				Name:   "Times Scolded Bot Person",
				Value:  fmt.Sprintf("%d Time(s)", userStats.BadBotCount),
				Inline: true,
			},
			{},
			{
				Name:   "Image(s) Requested",
				Value:  fmt.Sprintf("%d Image(s)", userStats.ImageCount),
				Inline: true,
			},
			{
				Name:   "Loot Box(s) Opened",
				Value:  fmt.Sprintf("%d Loot Box(s)", userStats.LootBoxCount),
				Inline: true,
			},
			{},
			{
				Name:   "Bonus Streak",
				Value:  fmt.Sprintf("%d Day(s)", user.BonusStreak),
				Inline: true,
			},
			{
				Name:   "Token Balance",
				Value:  fmt.Sprintf("%.2f Token(s)", user.ImageTokens),
				Inline: true,
			},
			{},
			// {
			// 	Name:  "Last Bonus",
			// 	Value: fmt.Sprintf("<t:%d:R>", userStats.LastBonus.Unix()),
			// },
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
