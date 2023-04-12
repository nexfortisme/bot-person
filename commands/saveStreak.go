package commands

import (
	"fmt"
	"main/persistance"
	"main/util"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SaveStreak(s *discordgo.Session, i *discordgo.InteractionCreate) {

	user, err := persistance.GetUserStats(i.Interaction.Member.User.ID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Something Went Wrong. Please start panicing.",
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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You don't need to save your streak.",
			},
		})
		return
	}

	// The time in seconds since the HoldStreatTimer was set
	diff := time.Now().Sub(user.HoldStreakTimer)

	// User timed out on saving streak. Reset the streak
	if diff.Seconds() >= 300 {

		// Setting random seed and generating a, value safe, token amount
		randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
		reward := randomizer.Intn(45) + 5
		rewardf64 := float64(reward) / 10.0
		finalReward := util.LowerFloatPrecision(rewardf64)

		// Updating User Record
		user.LastBonus = time.Now()
		user.ImageTokens += finalReward
		user.BonusStreak = 1
		user.HoldStreakTimer = time.Time{}

		resultString := fmt.Sprintf("Save Streak Timed Out. Streak Reset. \nCurrent Streak: 1 \nYou have been awarded %f tokens.", finalReward)

		result := persistance.UpdateUserStats(i.Interaction.Member.User.ID, user)
		if result {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: resultString,
				},
			})
			return
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something Went Wrong When Resetting Streak.",
				},
			})
			return
		}
	}

	if saveAction == "use" {

		if user.SaveStreakTokens <= 0 {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You don't have any Save tokens to use! Try again with the buy action if you really want to save your streak.",
				},
			})
			return
		}

		user.SaveStreakTokens -= 1

		user.BonusStreak++
		streak := user.BonusStreak

		if streak%10 == 0 && streak%100 != 0 && streak%50 != 0 {
			returnString = fmt.Sprintf("Congrats on keeping the streak alive. Current Streak: %d. Bonus Modifier: 2x", streak)
			modifier = 2
		} else if streak%25 == 0 && streak%50 != 0 && streak%100 != 0 {
			returnString = fmt.Sprintf("Great work on keeping the streak alive! Current Streak: %d. Bonus Modifier: 5x", streak)
			modifier = 5
		} else if streak%50 == 0 && streak%100 != 0 {
			returnString = fmt.Sprintf("Wow! That's a long time. Current Streak: %d. Bonus Modifier: 10x", streak)
			modifier = 10
		} else if streak%69 == 0 {
			returnString = fmt.Sprintf("Nice, Congratulations! Current Streak: %d. Bonus Modifier: 15x", streak)
			modifier = 15
		} else if streak%100 == 0 {
			returnString = fmt.Sprintf("Few people ever reach is this far, Congratulations! Current Streak: %d. Bonus Modifier: 50x", streak)
			modifier = 50
		} else {
			returnString = fmt.Sprintf("Current Bonus Streak: %d", streak)
		}

		// Setting random seed and generating a, value safe, token amount
		randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
		reward := randomizer.Intn(45) + 5
		reward *= modifier
		rewardf64 := float64(reward) / 10.0
		finalReward := util.LowerFloatPrecision(rewardf64)

		// Updating User Record
		user.LastBonus = time.Now()
		user.ImageTokens += finalReward
		user.HoldStreakTimer = time.Time{}

		finalString := fmt.Sprintf("STREAK SAVED! %s \n Congrats! You are awarded %.2f tokens and now have %d save tokens.", returnString, finalReward, user.SaveStreakTokens)

		result := persistance.UpdateUserStats(i.Interaction.Member.User.ID, user)
		if result {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: finalString,
				},
			})
			return
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something Went Wrong Saving Streak.",
				},
			})
			return
		}

	} else if saveAction == "buy" {

		// Using half the users tokens
		cost := user.ImageTokens / 2
		user.ImageTokens /= 2

		user.BonusStreak++
		streak := user.BonusStreak

		if streak%10 == 0 && streak%100 != 0 && streak%50 != 0 {
			returnString = fmt.Sprintf("Congrats on keeping the streak alive. Current Streak: %d. Bonus Modifier: 2x", streak)
			modifier = 2
		} else if streak%25 == 0 && streak%50 != 0 && streak%100 != 0 {
			returnString = fmt.Sprintf("Great work on keeping the streak alive! Current Streak: %d. Bonus Modifier: 5x", streak)
			modifier = 5
		} else if streak%50 == 0 && streak%100 != 0 {
			returnString = fmt.Sprintf("Wow! That's a long time. Current Streak: %d. Bonus Modifier: 10x", streak)
			modifier = 10
		} else if streak%69 == 0 {
			returnString = fmt.Sprintf("Nice, Congratulations! Current Streak: %d. Bonus Modifier: 15x", streak)
			modifier = 15
		} else if streak%100 == 0 {
			returnString = fmt.Sprintf("Few people ever reach is this far, Congratulations! Current Streak: %d. Bonus Modifier: 50x", streak)
			modifier = 50
		} else {
			returnString = fmt.Sprintf("Current Bonus Streak: %d", streak)
		}

		// Setting random seed and generating a, value safe, token amount
		randomizer := rand.New(rand.NewSource(time.Now().UnixMilli()))
		reward := randomizer.Intn(45) + 5
		reward *= modifier
		rewardf64 := float64(reward) / 10.0
		finalReward := util.LowerFloatPrecision(rewardf64)

		finalString := fmt.Sprintf("STREAK SAVED! Save Token Cost: %.2f. %s \n Congrats! You are awarded %.2f tokens and now have %d save tokens.", cost, returnString, finalReward, user.SaveStreakTokens)

		// Updating User Record
		user.LastBonus = time.Now()
		user.ImageTokens += finalReward
		user.HoldStreakTimer = time.Time{}

		result := persistance.UpdateUserStats(i.Interaction.Member.User.ID, user)
		if result {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: finalString,
				},
			})
			return
		} else {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Something Went Wrong Saving Streak.",
				},
			})
			return
		}

	}

}
