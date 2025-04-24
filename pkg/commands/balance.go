package commands

import (
	"fmt"
	"main/pkg/persistance"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type Balance struct{}

func (b *Balance) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "balance",
		Description: "Check your balance or the balance of another user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The user to check the balance of.",
				Required:    false,
			},
		},
	}
}

func (b *Balance) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
		balanceResponse = user.Username + " has " + fmt.Sprintf("%.2f", checkUser.ImageTokens) + " tokens."

		logging.LogEvent(eventType.COMMAND_BALANCE, i.Interaction.Member.User.ID, fmt.Sprintf("User has checked the balance of %s", user.ID), i.Interaction.GuildID)
	} else {

		user, _ := persistance.GetUser(i.Interaction.Member.User.ID)
		balanceResponse = "You have " + fmt.Sprintf("%.2f", user.ImageTokens) + " tokens."

		logging.LogEvent(eventType.COMMAND_BALANCE, i.Interaction.Member.User.ID, "User has checked their balance", i.Interaction.GuildID)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: balanceResponse,
		},
	})
}

func (b *Balance) HelpString() string {
	return "The `/balance` command allows you to see how many Bot Person tokens you currently have. You can also see the balance of others in the server with the `user` option. If you don't specify a user, it will default to you."
}

func (b *Balance) CommandCost() int {
	return 0
}
