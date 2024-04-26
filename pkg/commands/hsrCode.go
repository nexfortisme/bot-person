package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func HSRCode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.HSR_CODE, i.Interaction.Member.User.ID, "HSR Code Command Use", i.Interaction.GuildID)

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	codeResponse := "https://hsr.hoyoverse.com/gift?code="

	if option, ok := optionMap["prompt"]; ok {
		codeResponse += option.StringValue()
	}



	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: codeResponse,
		},
	})
}
