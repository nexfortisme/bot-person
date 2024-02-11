package commands

import (
	"fmt"
	"main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func Balance(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var tokenCount float64
	var balanceResponse string

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	if option, ok := optionMap["user"]; ok {
		user := option.UserValue(s)
		tokenCount = persistance.GetUserTokenCount(user.ID)
		balanceResponse = user.Username + " has " + fmt.Sprintf("%.2f", tokenCount) + " tokens."
	} else {
		tokenCount = persistance.GetUserTokenCount(i.Interaction.Member.User.ID)
		balanceResponse = "You have " + fmt.Sprintf("%.2f", tokenCount) + " tokens."
	}

	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: balanceResponse,
		},
	})
}