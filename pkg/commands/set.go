package commands

import (
	"fmt"
	persistance "main/pkg/persistance"
	attribute "main/pkg/persistance/eums"

	"github.com/bwmarrin/discordgo"
)

type Set struct{}

func (st *Set) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "set",
		Description: "Set a user attribute",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Name:        "attribute",
				Description: "The attribute to set",
				Type:        discordgo.ApplicationCommandOptionString,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "/Bot Pre-Prompt",
						Value: attribute.BOT_PREPROMPT.String(),
					},
				},
				Required: true,
			},
			{
				Name:        "value",
				Description: "The value to set the attribute to",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    false,
			},
		},
	}
}

func (st *Set) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	userId := i.Interaction.Member.User.ID
	attributeString := optionMap["attribute"].StringValue()

	attributeEnum := attribute.Attribute(attributeString)

	if optionMap["value"] == nil {
		// Get The Attribute
		attributeValue, err := persistance.GetUserAttribute(userId, attributeEnum)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error getting attribute",
				},
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(`The attribute: "%s" is: "%s"`, attributeEnum.String(), attributeValue),
			},
		})
	} else {

		value := optionMap["value"].StringValue()

		err := persistance.SetUserAttribute(userId, attributeEnum, value)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error setting attribute",
				},
			})
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(`The attribute: "%s" is now set to: "%s"`, attributeEnum.String(), value),
			},
		})
	}

}

func (st *Set) HelpString() string {
	return "The `/set` command allows you to set a user attribute. You can specify the attribute with the `attribute` option."
}

func (st *Set) CommandCost() int {
	return 0
}
