package commands

import (
	"fmt"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

var integerOptionMinValue = 0.1

type Send struct{}

func (sn *Send) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "send",
			Description: "Send tokens to another user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "recepient",
					Description: "The person you want to send tokens to.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of tokens you want to send.",
					MinValue:    &integerOptionMinValue,
					Required:    true,
				},
			},
	}
}

func (sn *Send) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender, _ := persistance.GetUser(i.Interaction.Member.User.ID)

	var transferrAmount float64

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["amount"]; ok {

		transferrAmount = option.FloatValue()

		if sender.ImageTokens < transferrAmount {

			logging.LogEvent(eventType.COMMAND_SEND, i.Interaction.Member.User.ID, "User does not have enough tokens to send", i.Interaction.GuildID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! You do not have the tokens needed to complete the transaction.",
				},
			})
			return
		}
	}

	if option, ok := optionMap["recepient"]; ok {
		recepient := option.UserValue(s)

		if i.Interaction.Member.User.ID == recepient.ID {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You cannot send tokens to yourself.",
				},
			})
			return
		}

		sendResponse := persistance.TransferBotPersonTokens(transferrAmount, i.Interaction.Member.User.ID, recepient.ID)

		sender, _ = persistance.GetUser(i.Interaction.Member.User.ID)

		if sendResponse {

			logging.LogEvent(eventType.COMMAND_SEND, i.Interaction.Member.User.ID, fmt.Sprintf("User Sent %f tokens to %s", transferrAmount, recepient.ID), i.Interaction.GuildID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Tokens were successfully sent. Your new balance is: " + fmt.Sprintf("%.2f", sender.ImageTokens),
				},
			})
			return
		} else {

			logging.LogEvent(eventType.COMMAND_SEND, i.Interaction.Member.User.ID, "Something Went Wrong.", i.Interaction.GuildID)

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

func (sn *Send) HelpString() string {
	return "The `/send` command allows you to send Bot Person tokens to other users in the server. You can specify the user to send to with the `user` option and the amount to send with the `amount` option."
}

func (sn *Send) CommandCost() int {
	return 0
}	
