package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Search struct{}

func (se *Search) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "search",
		Description: "Search the web for information",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The search query",
				Required:    true,
			},
		},
	}
}

func (se *Search) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var searchQuery string

	if option, ok := optionMap["prompt"]; ok {
		searchQuery = option.StringValue()

		logging.LogEvent(eventType.COMMAND_SEARCH, i.Interaction.Member.User.ID, searchQuery, i.Interaction.GuildID)

		placeholderBotResponse := "Searching for: " + searchQuery

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: placeholderBotResponse,
			},
		})

		perplexityResponse := external.GetPerplexityResponse("", searchQuery)

		if len(perplexityResponse.Choices) == 0 {
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Error getting response from Perplexity.",
			})
			return
		}

		response := perplexityResponse.Choices[0].Message.Content

		if perplexityResponse.Citations != nil {
			for index, citation := range perplexityResponse.Citations {
				replaceString := fmt.Sprintf("[%d]", index)
				replacementString := fmt.Sprintf("[[%d]](%s)", index, citation)
				response = strings.Replace(response, replaceString, replacementString, -1)
			}
		}

		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &response,
		})
	}

}

func (se *Search) HelpString() string {
	return "The `/search` command allows you to search the web for information. You can specify the search query with the `prompt` option."
}

func (se *Search) CommandCost() int {
	return 0
}
