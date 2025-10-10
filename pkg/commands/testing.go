package commands

import (
	"encoding/json"
	"fmt"
	"main/pkg/util"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Testing struct{}

func (im *Testing) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "testing",
		Description: fmt.Sprintf("I just wanna see the response"),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "foo_string",
				Description: "foo_string",
				Required:    true,
			},
		},
	}
}

func (im *Testing) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	userImageOptions := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	userImageOptionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(userImageOptions))
	for _, opt := range userImageOptions {
		userImageOptionMap[opt.Name] = opt
	}

	// Pulling the propt out of the optionsMap
	if option, ok := userImageOptionMap["foo_string"]; ok {
		fmt.Println(option.StringValue())

		currentTime := time.Now().Format("2006-01-02-15-04-05")

		data, _ := json.MarshalIndent(i, "", "  ")

		util.SaveResponseToFile(data, fmt.Sprintf("foo_string-%s.txt", currentTime))


		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Response saved to file.",
			},
		})
	}
}

func (im *Testing) HelpString() string {
	return fmt.Sprintf("foo_string")
}

func (im *Testing) CommandCost() int {
	return 0
}

