package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type HSRCode struct{}

func (h *HSRCode) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "hsr-code",
		Description: "Get the Honkai Star Rail gift code url from a code.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "The code to get the url for",
				Required:    true,
			},
		},
	}
}

func (h *HSRCode) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logging.LogEvent(eventType.HSR_CODE, i.Interaction.Member.User.ID, "HSR Code Command Use", i.Interaction.GuildID)

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	codeResponse := "https://hsr.hoyoverse.com/gift?code="

	if option, ok := optionMap["code"]; ok {
		codeResponse += option.StringValue()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: codeResponse,
		},
	})
}

func (h *HSRCode) HelpString() string {
	return "The `/hsr-code` command allows you to get the Honkai Star Rail gift code url from a code."
}

func (h *HSRCode) CommandCost() int {
	return 0
}
