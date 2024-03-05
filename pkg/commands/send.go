package commands

import (
	"fmt"
	persistance "main/pkg/persistance/services"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
)

func Send(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var transferrAmount float64
	senderStats, _ := persistance.GetUserStats(i.Interaction.Member.User.ID, s)


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

		if senderStats.Token_Balance < transferrAmount {

			logging.LogEvent(loggingType.COMMAND_SEND, "User attempted to send more tokens than they have", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

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

		recipientStats, _ := persistance.GetUserStats(recepient.ID, s)
		senderStats.Token_Balance = senderStats.Token_Balance - transferrAmount
		recipientStats.Token_Balance = recipientStats.Token_Balance + transferrAmount

		_, err := persistance.UpsertUserStats(senderStats)
		_, err = persistance.UpsertUserStats(recipientStats)

		if err == nil {

			logging.LogEvent(loggingType.COMMAND_SEND, "Tokens sent successfully", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Tokens were successfully sent. Your new balance is: " + fmt.Sprint(senderStats.Token_Balance),
				},
			})
			return
		} else {

			logging.LogEvent(loggingType.COMMAND_SEND, "Error sending tokens", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

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
