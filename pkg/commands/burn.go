package commands

import (
	"fmt"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Burn(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var burnAmount float64
	user, _ := persistance.GetUser(i.Interaction.Member.User.ID)

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

		if user.ImageTokens < burnAmount {

			logging.LogEvent(eventType.COMMAND_BURN, i.Interaction.Member.User.ID, "User does not have enough tokens to burn", i.Interaction.GuildID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Oops! You do not have the tokens needed to complete the transaction.",
				},
			})
			return
		} else {
			
			user.ImageTokens -= burnAmount
			persistance.UpdateUser(*user);

			logging.LogEvent(eventType.COMMAND_BURN, i.Interaction.Member.User.ID, fmt.Sprintf("User has burnt %f tokens", burnAmount), i.Interaction.GuildID)

			removeTokenResponse := fmt.Sprintf("%.2f tokens removed. New Balance: %.2f", burnAmount, user.ImageTokens)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: removeTokenResponse,
				},
			})
		}
	}

}