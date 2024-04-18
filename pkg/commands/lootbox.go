package commands

import (
	"fmt"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Lootbox(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

	logging.LogEvent(eventType.COMMAND_LOOTBOX, i.Interaction.Member.User.ID, fmt.Sprintf("User has purchased a lootbox with seed: %d and reward %f", lootboxSeed, lootboxReward), i.Interaction.GuildID)

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
