package commands

import (

	logging "main/pkg/logging/services"
	loggingType "main/pkg/logging/enums"

	persistance "main/pkg/persistance/services"

	"github.com/bwmarrin/discordgo"
)

func Burn(s *discordgo.Session, i *discordgo.InteractionCreate) {

	userStats, _ := persistance.GetUserStats(i.Interaction.Member.User.ID, s)
	userBalance := userStats.Token_Balance;

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["amount"]; ok {

		burnAmount = option.FloatValue()

		if userBalance < burnAmount {

			logging.LogEvent(loggingType.COMMAND_BURN, "User attempted to burn more tokens than they have", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! You do not have the tokens needed to complete the transaction.",
				},
			})
			return
		} else {

			userStats.Token_Balance = userBalance - burnAmount
			persistance.UpsertUserStats(userStats)
			userBalance = persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

			removeTokenResponse := fmt.Sprintf("%.2f tokens removed. New Balance: %.2f", burnAmount, senderBalance)

			logging.LogEvent(loggingType.COMMAND_BURN, removeTokenResponse, i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: removeTokenResponse,
				},
			})
		}
	}

}
