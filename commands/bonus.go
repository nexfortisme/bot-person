package commands

import (
	"fmt"
	"main/persistance"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Bonus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Bonus Reward: -1 or actual reward. -1 is just placeholder.
	// Return Message: "" or specialBonusRewardString. "" is default value and only returned if there is no modifier.
	// err: nil, wait error, or saveStreak error.
	bonusReward, returnMessage, err := persistance.GetUserReward(i.Interaction.Member.User.ID)
	var bonusReturnMessage string

	// If there is an error, then use the message from the error.
	if err != nil {
		bonusReturnMessage = err.Error()
	} else { // If there is no error, then use the message from the returnMessage.
		if returnMessage != "" { // If there is a returnMessage, then append the bonusReward to the returnMessage.
			bonusReturnMessage = fmt.Sprintf("%s \nCongrats! You are awarded %.2f tokens", returnMessage, bonusReward)
		} else if returnMessage == "" { // If there is no returnMessage, then just return the bonusReward.
			bonusReturnMessage = fmt.Sprintf("Congrats! You are awarded %.2f tokens", bonusReward)
		}
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: bonusReturnMessage,
		},
	})

	// Cleaning up the bonus message if the user is on cooldown or missed their bonus window.
	if err != nil {
		time.Sleep(time.Second * 15)
		_ = s.InteractionResponseDelete(i.Interaction)
	}

}
