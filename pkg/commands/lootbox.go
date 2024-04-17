package commands

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Lootbox(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	lootboxReward, lootboxSeed, err := persistance.BuyLootbox(i.Interaction.Member.User.ID)
	var lootboxReturnMessage string

	if err != nil {
		lootboxReturnMessage = err.Error()
	} else {

		// TODO - Refactor this so a change in rates doesn't break the command
		if lootboxReward == 3.63 {
			lootboxReturnMessage = fmt.Sprintf("%s You purchased a lootbox with the seed: %d and it contained %.2f tokens", util.GetOofResponse(), lootboxSeed, lootboxReward)
		} else if lootboxReward == 8 {
			lootboxReturnMessage = fmt.Sprintf("You purchased a lootbox with the seed: %d and it contained %f tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 15 {
			lootboxReturnMessage = fmt.Sprintf("Congrats! You purchased a lootbox with the seed: %d and it contained %f tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 25 {
			lootboxReturnMessage = fmt.Sprintf("Woah! You purchased a lootbox with the seed: %d and it contained %f tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 50 {
			lootboxReturnMessage = fmt.Sprintf("Stop Hacking. You purchased a lootbox with the seed: %d and it contained %f tokens", lootboxSeed, lootboxReward)
		}

	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: lootboxReturnMessage,
		},
	})

	// Cleaning up the bonus message if the user is on cooldown
	if err != nil {
		time.Sleep(time.Second * 15)
		s.InteractionResponseDelete(i.Interaction)
	}

}
