package handlers

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SaveStreakButton(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Parse CustomID to retrieve the stored user ID
	parts := strings.Split(i.MessageComponentData().CustomID, ":")
	if len(parts) < 2 {
		fmt.Println("Invalid CustomID format")
		return
	}

	originalUserID := parts[1]
	clickingUserID := i.Member.User.ID

	if originalUserID != clickingUserID {
		// Respond with a message that the user does not have permission to click this button
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not allowed to interact with this button.",
				Flags:   discordgo.MessageFlagsEphemeral, // Make the response only visible to the user
			},
		})
		return
	}

	// Getting User
	user, _ := persistance.GetUser(originalUserID);

	var saveStreakCost float64
	var saveStreakMessage string

	// Calculating cost and creating save streak string
	saveStreakCost = user.UserStats.ImageTokens * 0.1
	saveStreakMessage = fmt.Sprintf("You have saved your streak! It Cost %.2f tokens", saveStreakCost)

	// Removing tokens from user
	persistance.RemoveBotPersonTokens(saveStreakCost, originalUserID)

	// Refetching stats
	user, _ = persistance.GetUser(i.Interaction.Member.User.ID)

	// Updating the streak
	user.UserStats.BonusStreak++

	//Getting return string and modifier
	_, modifier := util.GetStreakStringAndModifier(user.UserStats.BonusStreak)

	// Getting Final Bonus Reward
	finalReward := util.GetUserBonus(5, 50, modifier)

	// Updating User Record
	user.UserStats.LastBonus = time.Now()
	user.UserStats.ImageTokens += finalReward

	// Updating User Stats
	persistance.UpdateUser(*user)

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
				Value: fmt.Sprintf("%d days", user.UserStats.BonusStreak),
			},
			{
				Name:  "Current Balance",
				Value: fmt.Sprintf("%.2f tokens", user.UserStats.ImageTokens),
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
	timerManager.ExecTimerFunction("save_streak_exp" + i.Interaction.Member.User.ID)

	// s.InteractionResponseDelete(i.Interaction)

}
