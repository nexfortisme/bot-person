package commands

import (
	"main/pkg/external"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
)

func Bot(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
		botResponseString = ParseSlashCommand(s, option.StringValue())

		logging.LogEvent(loggingType.COMMAND_BOT, botResponseString, i.Member.User.ID, i.GuildID, s)

		// Updating the initial message with the response from the OpenAI API
		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &botResponseString,
		})
		if err != nil {

			logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "Error editing interaction response for /bot", i.Member.User.ID, i.GuildID, s)

			// Not 100% sure this is the approach I want to take with handling errors from the API
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
			return
		}
	}
}

func ParseSlashCommand(s *discordgo.Session, prompt string) string {
	respTxt := external.GetOpenAIResponse(prompt)
	respTxt = "Request: " + prompt + " " + respTxt
	if len("Request: "+prompt+" "+respTxt) > 2000 {
		respTxt = respTxt[:1997] + "..."
	}
	return respTxt
}
