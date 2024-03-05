package commands

import (
	// "main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func Store(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// user, err := persistance.GetUserStats(i.Interaction.Member.User.ID)
	// if err != nil {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Something Went Wrong. Please start panicing.",
	// 		},
	// 	})
	// 	return
	// }

	// var purchaseItem string

	// // Access options in the order provided by the user.
	// options := i.ApplicationCommandData().Options

	// // Or convert the slice into a map
	// optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	// for _, opt := range options {
	// 	optionMap[opt.Name] = opt
	// }

	// // Checking to see that the user has the number of tokens needed to send
	// if option, ok := optionMap["item"]; ok {
	// 	purchaseItem = option.StringValue()
	// }

	// if purchaseItem == "help" {
	// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 		Data: &discordgo.InteractionResponseData{
	// 			Content: "Current Stock: \nSave Streak Token: 50 Tokens",
	// 		},
	// 	})
	// 	return
	// }

	// if purchaseItem == "save-streak-token" {

	// 	if user.ImageTokens < 50 {
	// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 			Data: &discordgo.InteractionResponseData{
	// 				Content: "You don't have enough tokens to buy this item.",
	// 			},
	// 		})
	// 		return
	// 	}

	// 	user.ImageTokens -= 50
	// 	user.SaveStreakTokens += 1

	// 	result := persistance.UpdateUserStats(i.Interaction.Member.User.ID, user)
	// 	if result {
	// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 			Data: &discordgo.InteractionResponseData{
	// 				Content: "Save Streak Token Purchased.",
	// 			},
	// 		})
	// 		return
	// 	} else {
	// 		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
	// 			Data: &discordgo.InteractionResponseData{
	// 				Content: "Something Went Wrong When Purchasing Item.",
	// 			},
	// 		})
	// 		return
	// 	}

	// }

}
