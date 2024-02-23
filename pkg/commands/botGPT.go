package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"

	"github.com/bwmarrin/discordgo"
)

func BotGPT(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var botResponseString string

	// Pulling the propt out of the optionsMap
	if option, ok := optionMap["prompt"]; ok {

		// Generating the response
		placeholderBotResponse := "Thinking about: " + option.StringValue()

		// Immediately responding in the 3 second window before the interaciton times out
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: placeholderBotResponse,
			},
		})

		// Going out to make the OpenAI call to get the proper response
		// botResponseString = ParseGPTSlashCommand(s, option.StringValue(
		// Check if the response will be too long and truncate if necessary
		prompt := option.StringValue()
		botResponseString = ParseGPTSlashCommand(s, prompt)
		if len("Request: "+prompt+" ")+len(botResponseString) > 2000 {
			truncatedLength := 2000 - len(prompt) - len("...") // account for ellipsis
			if truncatedLength > 0 {
				botResponseString = botResponseString[:truncatedLength] + "..."
			} else {
				botResponseString = "Response too long to display."
			}
		}

		// Incrementint interaciton counter
		persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *i.Interaction.Member.User)

		// Updating the initial message with the response from the OpenAI API
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &botResponseString,
		})
		if err != nil {

			fmt.Println("Error editing interaction response:", err)

			// Not 100% sure this is the approach I want to take with handling errors from the API
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
			return
		}
	}
}

func ParseGPTSlashCommand(s *discordgo.Session, prompt string) string {
	respTxt := external.GetOpenAIGPTResponse(prompt)
	respTxt = "Request: " + prompt + " " + respTxt
	return respTxt
}
