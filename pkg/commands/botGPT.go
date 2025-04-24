package commands

import (
	"fmt"
	"main/pkg/external"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type BotGPT struct {}

func (b *BotGPT) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "bot-gpt",
		Description: "Interact with OpenAI's GPT-4 API and see what out future AI overlords have to say.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "prompt",
				Description: "The actual prompt that the bot will ponder on.",
				Required: true,
			},
		},
	}
}

func (b *BotGPT) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

		logging.LogEvent(eventType.COMMAND_BOT_GPT, i.Interaction.Member.User.ID, option.StringValue(), i.Interaction.GuildID)

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
			truncatedLength := 2000 - len("Request: "+prompt+" ") - len("...") // account for ellipsis
			if truncatedLength > 0 {
				botResponseString = botResponseString[:truncatedLength] + "..."
			} else {
				botResponseString = "Response too long to display."
			}
		}

		logging.LogEvent(eventType.EXTERNAL_GPT_RESPONSE, i.Interaction.Member.User.ID, botResponseString, i.Interaction.GuildID)

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

func (b *BotGPT) HelpString() string {
	return "The `/bot-gpt` command allows you to prompt OpenAI's GPT-3 or GPT-4 chat model. You can ask it whatever as part of the `prompt` and once it generates a response, it will update the message with what came back. This is slower than the `/bot` command due to the chat model being more complex."
}

func (b *BotGPT) CommandCost() int {
	return 0
}

func ParseGPTSlashCommand(s *discordgo.Session, prompt string) string {
	respTxt := external.GetOpenAIGPTResponse(prompt)
	respTxt = "Request: " + prompt + " " + respTxt
	return respTxt
}
