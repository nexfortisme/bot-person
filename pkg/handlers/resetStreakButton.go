package handlers

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func ResetStreakButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Getting User Stats
	userStats, userStatsError := persistance.GetUserStats(i.Interaction.Member.User.ID)
	if userStatsError != nil {
		return
	}

	// Resetting the streak
	userStats.BonusStreak = 1

	//Getting return string and modifier
	_, modifier := util.GetStreakStringAndModifier(userStats.BonusStreak)

	// Getting Final Bonus Reward
	finalReward := util.GetUserBonus(5, 50, modifier)

	// Updating User Record
	userStats.LastBonus = time.Now()
	userStats.ImageTokens += finalReward

	// Updating User Stats
	persistance.UpdateUserStats(i.Interaction.Member.User.ID, userStats)

	resetStreakEmbed := &discordgo.MessageEmbed{
		Title:       "Reset Streak",
		Description: "You have reset your streak!",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Bonus Award",
				Value: fmt.Sprintf("%.2f tokens", finalReward),
			},
			{
				Name:  "Current Streak",
				Value: fmt.Sprintf("%d days", userStats.BonusStreak),
			},
			{
				Name:  "Current Balance",
				Value: fmt.Sprintf("%.2f tokens", userStats.ImageTokens),
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{resetStreakEmbed},
		},
	})

	timerManager := util.GetInstance()
	timerManager.StopTimer("save_streak_exp" + i.Interaction.Member.User.ID)

}
