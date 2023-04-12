package commands

import (
	"fmt"
	"main/persistance"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Bonus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	bonusReward, returnMessage, err := persistance.GetUserReward(i.Interaction.Member.User.ID)
	var bonusReturnMessage string

	if err != nil {
		bonusReturnMessage = err.Error()
	} else {

		if bonusReward == -1 {
			bonusReturnMessage = returnMessage
		}

		if returnMessage != "" && bonusReward != -1 {
			bonusReturnMessage = fmt.Sprintf("%s \nCongrats! You are awarded %.2f tokens", returnMessage, bonusReward)
		} else if returnMessage == "" && bonusReward != -1 {
			bonusReturnMessage = fmt.Sprintf("Congrats! You are awarded %.2f tokens", bonusReward)
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: bonusReturnMessage,
		},
	})

	// Cleaning up the bonus message if the user is on cooldown
	if err != nil {
		time.Sleep(time.Second * 15)
		s.InteractionResponseDelete(i.Interaction)
	}

}
