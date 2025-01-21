package commands

import (
	"fmt"
	"main/pkg/persistance"
	persistanceEnums "main/pkg/persistance/eums"
	"main/pkg/util"
	"time"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Bonus(s *discordgo.Session, i *discordgo.InteractionCreate) {

	logging.LogEvent(eventType.COMMAND_BONUS, i.Interaction.Member.User.ID, "User has checked their bonus", i.Interaction.GuildID)

	// Bonus Reward: -1 or actual reward. -1 is just placeholder.
	// Return Message: "" or specialBonusRewardString. "" is default value and only returned if there is no modifier.
	// err: nil, wait error, or saveStreak error.
	bonusReward, rewardStatus, err := persistance.GetUserReward(i.Interaction.Member.User.ID)

	user, _ := persistance.GetUser(i.Interaction.Member.User.ID)

	// Variables for the embed message.
	var flavorText string
	var embedFields []*discordgo.MessageEmbedField
	var bonusEmbed *discordgo.MessageEmbed
	var disableSaveStreak = false

	var TEM_MINUTES = 10 * time.Minute

	if rewardStatus == persistanceEnums.TOO_EARLY {
		embedFields = []*discordgo.MessageEmbedField{
			{
				Name: "You already collected your bonus today!",
			},
			{
				Name:  "Next Bonus In",
				Value: err.Error(),
			},
		}
	} else if rewardStatus == persistanceEnums.MISSED {
		embedFields = []*discordgo.MessageEmbedField{
			{
				Name:  "Current Streak",
				Value: fmt.Sprintf("%d days", user.BonusStreak),
			},
			{
				Name:  "Streak Missed",
				Value: "Click the save streak button to save your streak and get your bonus. Save streak costs 10% of your current balance. Click 'Save Streak' in the next 10 minutes to save your streak, or click 'Reset Streak' to reset your streak at no cost.",
			},
			{
				Name:  "Current Balance",
				Value: fmt.Sprintf("%.2f tokens", user.ImageTokens),
			},
		}
	} else {
		flavorText, _ = util.GetStreakStringAndModifier(user.BonusStreak)
		embedFields = []*discordgo.MessageEmbedField{
			{
				Name:  "Bonus Award",
				Value: fmt.Sprintf("%.2f tokens", bonusReward),
			},
			{
				Name:  "Current Streak",
				Value: fmt.Sprintf("%d days", user.BonusStreak),
			},
			{
				Name:  "Current Balance",
				Value: fmt.Sprintf("%.2f tokens", user.ImageTokens),
			},
		}
	}

	var embedTitle string
	if rewardStatus == persistanceEnums.TOO_EARLY {
		embedTitle = "Too Early"
	} else if rewardStatus == persistanceEnums.MISSED {
		embedTitle = "Missed Bonus"
	} else {
		embedTitle = "Bonus Reward!"
	}

	var embedColor int
	if rewardStatus == persistanceEnums.TOO_EARLY {
		embedColor = 0xFFFF00 // Yellow
	} else if rewardStatus == persistanceEnums.MISSED {
		embedColor = 0xFF0000 // Red
	} else {
		embedColor = 0x00FF00 // Green
	}

	bonusEmbed = &discordgo.MessageEmbed{
		Title:       embedTitle,
		Description: flavorText,
		Fields:      embedFields,
		Color:       embedColor,
	}

	var buttonLabel string
	if rewardStatus == persistanceEnums.MISSED {
		buttonLabel = fmt.Sprintf("Save Streak (%.2f tokens)", user.ImageTokens/10)
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

	if rewardStatus == persistanceEnums.TOO_EARLY || rewardStatus == persistanceEnums.AVAILABLE {
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
