package commands

import (
	"fmt"
	"main/persistance"

	"github.com/bwmarrin/discordgo"
)

func Send(s *discordgo.Session, i *discordgo.InteractionCreate) {

	senderBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

	var transferAmount float64

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["amount"]; ok {

		transferAmount = option.FloatValue()

		if senderBalance < transferAmount {

			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! You do not have the tokens needed to complete the transaction.",
				},
			})
			return
		}
	}

	if option, ok := optionMap["recipient"]; ok {
		recipient := option.UserValue(s)

		if i.Interaction.Member.User.ID == recipient.ID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot send tokens to yourself.",
				},
			})
			return
		}

		sendResponse := persistance.TransferBotPersonTokens(transferAmount, i.Interaction.Member.User.ID, recipient.ID)

		newBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

		if sendResponse {

			// TODO - Switch to use BPSystemInteraction
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Tokens were successfully sent. Your new balance is: " + fmt.Sprint(newBalance),
				},
			})
			return
		} else {

			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! Something went wrong. Tokens were not sent.",
				},
			})
			return
		}
	}
}
