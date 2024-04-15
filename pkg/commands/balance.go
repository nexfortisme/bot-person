package commands

import (
	"fmt"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	persistance "main/pkg/persistance/services"

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
		queryUserStats, _ := persistance.GetUserByDiscordUserId(user.ID, s)

		balanceResponse = user.Username + " has " + fmt.Sprintf("%.2f", queryUserStats.UserStats.ImageTokens) + " tokens."
	} else {
		queryUserStats, _ := persistance.GetUserByDiscordUserId(i.Interaction.Member.User.ID, s)

		balanceResponse = "You have " + fmt.Sprintf("%.2f", queryUserStats.UserStats.ImageTokens) + " tokens."
	}

	logging.LogEvent(loggingType.COMMAND_BALANCE, balanceResponse, i.Member.User.ID, i.GuildID, s)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: balanceResponse,
		},
	})
}
