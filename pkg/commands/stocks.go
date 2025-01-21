package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"

	"github.com/bwmarrin/discordgo"
)

func Stocks(s *discordgo.Session, i *discordgo.InteractionCreate) {

	user, _ := persistance.GetUser(i.Interaction.Member.User.ID);

	var stockAction string
	var stockTicker string
	var purchaseAmount float64

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["action"]; ok {
		stockAction = option.StringValue()
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["stock"]; ok {
		stockTicker = option.StringValue()
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["quantity"]; ok {
		purchaseAmount = option.FloatValue()
	}

	// Getting the current price of the stock
	currentPrice, err := external.GetStockData(stockTicker)
	currentPriceF64 := util.LowerFloatPrecision(float64(currentPrice))

	// error with either the stock ticker or the API
	if err != nil {
		retString := fmt.Sprintf("Error getting stock data for ticker %s. Please try again.", stockTicker)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: retString,
			},
		})
		return
	}

	if stockAction == "buy" {

		purchasePrice := currentPriceF64 * purchaseAmount

		if user.ImageTokens < purchasePrice {

			retString := fmt.Sprintf("You don't have the tokens needed to purchase %.2f shares of %s. Please try again with a lower amount.", purchaseAmount, stockTicker)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: retString,
				},
			})
			return

		} else {

			// persistance.AddStock(i.Interaction.Member.User.ID, stockTicker, purchaseAmount)

			persistance.RemoveBotPersonTokens(purchasePrice, i.Interaction.Member.User.ID)

			user, _ = persistance.GetUser(i.Interaction.Member.User.ID)

			retString := fmt.Sprintf("You have purchased %f shares of %s for %.2f tokens. Your new balance is %.2f tokens.", purchaseAmount, stockTicker, purchasePrice, user.ImageTokens)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: retString,
				},
			})
			return
		}

	} else {

		// userStock, err := persistance.GetUserStock(i.Interaction.Member.User.ID, stockTicker)

		if err != nil {
			retString := fmt.Sprintf("You do not have any shares of %s. Please try again.", stockTicker)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: retString,
				},
			})
			return
		}

		// if userStock.StockCount < purchaseAmount {

		// 	retString := fmt.Sprintf("You do not have enough shares of %s to sell %.2f shares. Please try again with a lower amount.", stockTicker, purchaseAmount)

		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
		// 		Data: &discordgo.InteractionResponseData{
		// 			Content: retString,
		// 		},
		// 	})
		// 	return

		// } else {

		// 	sellPrice := currentPriceF64 * purchaseAmount

		// 	// persistance.RemoveStock(i.Interaction.Member.User.ID, stockTicker, purchaseAmount)

		// 	persistance.AddBotPersonTokens(sellPrice, i.Interaction.Member.User.ID)

		// 	user, _ = persistance.GetUser(i.Interaction.Member.User.ID)

		// 	retString := fmt.Sprintf("You have sold %f shares of %s for %.2f tokens. Your new balance is %.2f tokens.", purchaseAmount, stockTicker, sellPrice, user.ImageTokens)

		// 	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
		// 		Data: &discordgo.InteractionResponseData{
		// 			Content: retString,
		// 		},
		// 	})
		// 	return

		// }

	}
}