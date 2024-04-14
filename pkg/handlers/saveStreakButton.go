package handlers

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SaveStreakButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Getting User Stats
	userStats, userStatsError := persistance.GetUserStats(i.Interaction.Member.User.ID)
	if userStatsError != nil {
		return
	}

	var saveStreakCost float64
	var saveStreakMessage string

	// Calculating cost and creating save streak string
	saveStreakCost = userStats.ImageTokens * 0.1
	saveStreakMessage = fmt.Sprintf("You have saved your streak! It Cost %.2f tokens", saveStreakCost)

	// Removing tokens from user
	persistance.RemoveUserTokens(i.Interaction.Member.User.ID, saveStreakCost)

	// Refetching stats
	userStats, userStatsError = persistance.GetUserStats(i.Interaction.Member.User.ID)
	if userStatsError != nil {
		return
	}

	// Updating the streak
	userStats.BonusStreak++

	//Getting return string and modifier
	_, modifier := util.GetStreakStringAndModifier(userStats.BonusStreak)

	// Getting Final Bonus Reward
	finalReward := util.GetUserBonus(5, 50, modifier)

	// Updating User Record
	userStats.LastBonus = time.Now()
	userStats.ImageTokens += finalReward

	// Updating User Stats
	persistance.UpdateUserStats(i.Interaction.Member.User.ID, userStats)

	saveStreakEmbed := &discordgo.MessageEmbed{
		Title:       "Streak Saved!",
		Description: saveStreakMessage,
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

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{saveStreakEmbed},
		},
	})

	timerManager := util.GetInstance()
	timerManager.StopTimer("save_streak_exp" + i.Interaction.Member.User.ID)

	// s.InteractionResponseDelete(i.Interaction)

}
