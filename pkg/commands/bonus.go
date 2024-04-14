package commands

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Bonus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Bonus Reward: -1 or actual reward. -1 is just placeholder.
	// Return Message: "" or specialBonusRewardString. "" is default value and only returned if there is no modifier.
	// err: nil, wait error, or saveStreak error.
	bonusReward, _, err := persistance.GetUserReward(i.Interaction.Member.User.ID)

	// Variables for the embed message.
	// var bonusReturnMessage string
	var userStreak int
	var userBalance float64
	var flavorText string
	var embedFields []*discordgo.MessageEmbedField
	var bonusEmbed *discordgo.MessageEmbed
	var disableSaveStreak = false

	var TEM_MINUTES = 10 * time.Minute

	//Getting user stat information
	userBalance = persistance.GetUserTokenCount(i.Interaction.Member.User.ID)
	userStats, userStatsError := persistance.GetUserStats(i.Interaction.Member.User.ID)

	// If there is an error, then set the userStreak to 0.
	if userStatsError != nil {
		userStreak = 0
	} else {
		userStreak = userStats.BonusStreak
	}

	// // If there is an error, then use the message from the error.
	// if err != nil {
	// 	bonusReturnMessage = err.Error()
	// } else { // If there is no error, then use the message from the returnMessage.
	// 	if returnMessage != "" { // If there is a returnMessage, then append the bonusReward to the returnMessage.
	// 		bonusReturnMessage = fmt.Sprintf("%s \nCongrats! You are awarded %.2f tokens", returnMessage, bonusReward)
	// 	} else if returnMessage == "" { // If there is no returnMessage, then just return the bonusReward.
	// 		bonusReturnMessage = fmt.Sprintf("Congrats! You are awarded %.2f tokens", bonusReward)
	// 	}
	// }

	if err != nil {
		if bonusReward == -1 {
			embedFields = []*discordgo.MessageEmbedField{
				{
					Name: "You already collected your bonus today!",
				},
				{
					Name:  "Next Bonus In",
					Value: err.Error(),
				},
			}
		} else {

			// flavorText = "You Missed Your Bonus!"
			embedFields = []*discordgo.MessageEmbedField{
				{
					Name:  "Current Streak",
					Value: fmt.Sprintf("%d days", userStreak),
				},
				{
					Name:  "Streak Missed",
					Value: "Click the save streak button to save your streak and get your bonus. Save streak costs 10% of your current balance. Click 'Save Streak' in the next 10 minutes to save your streak, or click 'Reset Streak' to reset your streak at no cost.",
				},
				{
					Name:  "Current Balance",
					Value: fmt.Sprintf("%.2f tokens", userBalance),
				},
			}
		}
	} else {
		flavorText, _ = util.GetStreakStringAndModifier(userStreak)
		embedFields = []*discordgo.MessageEmbedField{
			{
				Name:  "Bonus Award",
				Value: fmt.Sprintf("%.2f tokens", bonusReward),
			},
			{
				Name:  "Current Streak",
				Value: fmt.Sprintf("%d days", userStreak),
			},
			{
				Name:  "Current Balance",
				Value: fmt.Sprintf("%.2f tokens", userBalance),
			},
		}
	}

	var embedTitle string
	if err != nil {
		embedTitle = "Missed Bonus"
	} else {
		embedTitle = "Bonus Reward!"
	}

	var embedColor int
	if err != nil {
		embedColor = 0xff0000
	} else {
		embedColor = 0x0000FF
	}

	bonusEmbed = &discordgo.MessageEmbed{
		Title:       embedTitle,
		Description: flavorText,
		Fields:      embedFields,
		Color:       embedColor,
	}

	var buttonLabel string
	if err != nil {
		buttonLabel = fmt.Sprintf("Save Streak (%.2f tokens)", userBalance/10)
	} else {
		buttonLabel = "Save Streak"
	}

	saveStreakButton := discordgo.Button{
		Label:    buttonLabel,
		Emoji:    &discordgo.ComponentEmoji{Name: "🍻"}, // This is needed to get the button to work
		Style:    discordgo.PrimaryButton,
		CustomID: "save_streak_button:" + i.Interaction.Member.User.ID,
		Disabled: err == nil || disableSaveStreak,
	}

	resetStreakButton := discordgo.Button{
		Label:    "Reset Streak",
		Emoji:    &discordgo.ComponentEmoji{Name: "🔥"}, // This is needed to get the button to work
		Style:    discordgo.DangerButton,
		CustomID: "reset_streak_button:" + i.Interaction.Member.User.ID,
		Disabled: err == nil,
	}

	if err != nil && bonusReward == -1 {
		saveStreakButton.Disabled = true
		resetStreakButton.Disabled = true
	}
	actionRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{saveStreakButton, resetStreakButton},
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{bonusEmbed},
			Components: []discordgo.MessageComponent{actionRow},
		},
	})

	if err != nil {
		timeString := "save_streak_exp" + i.Interaction.Member.User.ID
		timerManager := util.GetInstance()

		timerManager.SetTimer(timeString, TEM_MINUTES, func() {
			saveStreakButton.Disabled = true
			resetStreakButton.Disabled = true
			actionRow = discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{saveStreakButton, resetStreakButton},
			}

			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds:     &[]*discordgo.MessageEmbed{bonusEmbed},
				Components: &[]discordgo.MessageComponent{actionRow},
			})
			if err != nil {
				fmt.Println("error editing message", err)
			}
		})
	}
}
