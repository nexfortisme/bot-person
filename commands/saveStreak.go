package commands

import (
	"fmt"
	"main/persistance"
	"main/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SaveStreak(s *discordgo.Session, i *discordgo.InteractionCreate) {

	user, err := persistance.GetUserStats(i.Interaction.Member.User.ID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something Went Wrong. Please start panicking.",
			},
		})
		return
	}

	var saveAction string
	var returnString string
	var modifier int = 1

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["action"]; ok {
		saveAction = option.StringValue()
	}

	// If the holdStreakTimer is not set, then it bails out
	if (user.HoldStreakTimer == time.Time{}) {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't need to save your streak.",
			},
		})
		return
	}

	// The time in seconds since the HoldStreakTimer was set
	diff := time.Now().Sub(user.HoldStreakTimer)

	// User timed out on saving streak. Reset the streak
	if diff.Seconds() >= 300 {

		// Updating User Record
		user.LastBonus = time.Now()
		user.HoldStreakTimer = time.Time{}
		user.BonusStreak = 1

		// Calculating the reward
		finalReward := util.GetUserBonus(5, 50, modifier)

		// Updating user tokens
		user.ImageTokens += finalReward

		resultString := fmt.Sprintf("Save Streak Timed Out. Streak Reset. \nCurrent Streak: 1 \nYou have been awarded %.2f tokens.", finalReward)

		// if something went wrong when updating the user stats, then return an error message
		if !persistance.UpdateUserStats(i.Interaction.Member.User.ID, user) {
			returnString = "Something Went Wrong When Resetting Streak."
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: resultString,
			},
		})
		return
	}

	// The user opted to use a save token
	if saveAction == "use" {

		// User doesn't have any save tokens
		if user.SaveStreakTokens <= 0 {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You don't have any Save tokens to use! Try again with the buy action if you really want to save your streak.",
				},
			})
			return
		}

		user.SaveStreakTokens--
		user.BonusStreak++

		user.HoldStreakTimer = time.Time{}
		user.LastBonus = time.Now()

		returnString, modifier = util.GetStreakStringAndModifier(user.BonusStreak)

		finalReward := util.GetUserBonus(5, 50, modifier)

		// Updating User Tokens
		user.ImageTokens += finalReward

		finalString := fmt.Sprintf("STREAK SAVED! %s \n Congrats! You are awarded %.2f tokens and now have %d save tokens.", returnString, finalReward, user.SaveStreakTokens)

		if !persistance.UpdateUserStats(i.Interaction.Member.User.ID, user) {
			finalString = "Something Went Wrong Saving Streak."
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: finalString,
			},
		})
		return

	} else if saveAction == "buy" {

		// Using half the users tokens
		cost := user.ImageTokens / 2
		user.ImageTokens /= 2

		user.LastBonus = time.Now()
		user.HoldStreakTimer = time.Time{}

		user.BonusStreak++

		returnString, modifier = util.GetStreakStringAndModifier(user.BonusStreak)

		finalReward := util.GetUserBonus(5, 50, modifier)

		finalString := fmt.Sprintf("STREAK SAVED! Save Token Cost: %.2f. %s \n Congrats! You are awarded %.2f tokens.", cost, returnString, finalReward)

		// Updating User Tokens
		user.ImageTokens += finalReward

		if !persistance.UpdateUserStats(i.Interaction.Member.User.ID, user) {
			finalString = "Something Went Wrong Saving Streak."
		}

		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: finalString,
			},
		})
		return

	}

}
