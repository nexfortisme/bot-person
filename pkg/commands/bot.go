package commands

import (
	"fmt"
	"main/pkg/external"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {}

func (b *Bot) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "bot",
		Description: "A command to ask the bot for a response from their infinite wisdom.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "prompt",
				Description: "The actual prompt that the bot will ponder on.",
			},
		},
	}
}

func (b *Bot) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

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

		logging.LogEvent(eventType.COMMAND_BOT, i.Interaction.Member.User.ID, option.StringValue(), i.Interaction.GuildID)

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
		botResponseString = parseSlashCommand(option.StringValue(), i.Interaction.Member.User.ID)

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

func (b *Bot) HelpString() string {
	return "A command to ask the bot for a response from their infinite wisdom."
}

func parseSlashCommand(prompt string, userId string) string {
	respTxt := external.GetOpenAIResponse(prompt, userId)
	respTxt = "Request: " + prompt + " " + respTxt
	return respTxt
}