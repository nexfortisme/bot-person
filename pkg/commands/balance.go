package commands

import (
	"fmt"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Balance(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

		checkUser, _ := persistance.GetUser(user.ID)
		balanceResponse = user.Username + " has " + fmt.Sprintf("%.2f", checkUser.UserStats.ImageTokens) + " tokens."

		logging.LogEvent(eventType.COMMAND_BALANCE, i.Interaction.Member.User.ID, fmt.Sprintf("User has checked the balance of %s", user.ID), i.Interaction.GuildID)
	} else {
		
		user, _ := persistance.GetUser(i.Interaction.Member.User.ID)
		balanceResponse = "You have " + fmt.Sprintf("%.2f", user.UserStats.ImageTokens) + " tokens."

		logging.LogEvent(eventType.COMMAND_BALANCE, i.Interaction.Member.User.ID, "User has checked their balance", i.Interaction.GuildID)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: balanceResponse,
		},
	})
}