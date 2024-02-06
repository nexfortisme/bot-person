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
		if lootboxReward == 2 {
			lootboxReturnMessage = fmt.Sprintf("%s You purchased a lootbox with the seed: %d and it contained %d tokens", util.GetOofResponse(), lootboxSeed, lootboxReward)
		} else if lootboxReward == 5 {
			lootboxReturnMessage = fmt.Sprintf("You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 20 {
			lootboxReturnMessage = fmt.Sprintf("Congrats! You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 100 {
			lootboxReturnMessage = fmt.Sprintf("Woah! You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
		} else if lootboxReward == 500 {
			lootboxReturnMessage = fmt.Sprintf("Stop Hacking. You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
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
