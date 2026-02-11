package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type BotGPT struct{}

func (b *BotGPT) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "bot-gpt",
		Description: "Interact with OpenAI's GPT-4 API and see what out future AI overlords have to say.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "The actual prompt that the bot will ponder on.",
				Required:    true,
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
	var assistantResponse string

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
		botResponseString, assistantResponse = ParseGPTSlashCommand(prompt)

		if len(botResponseString) > 2000 {

			fileObj := util.HandleTooLongResponse(botResponseString)

			// Updating the initial message with the response from the OpenAI API
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Files: []*discordgo.File{fileObj},
			})

			if err != nil {
				// Not 100% sure this is the approach I want to take with handling errors from the API
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went oopsie with sending the file.",
				})
				return
			}
		} else {
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

		// if len("Request: "+prompt+" ")+len(botResponseString) > 2000 {
		// 	truncatedLength := 2000 - len("Request: "+prompt+" ") - len("...") // account for ellipsis
		// 	if truncatedLength > 0 {
		// 		botResponseString = botResponseString[:truncatedLength] + "..."
		// 	} else {
		// 		botResponseString = "Response too long to display."
		// 	}
		// }

		logging.LogEvent(eventType.EXTERNAL_GPT_RESPONSE, i.Interaction.Member.User.ID, assistantResponse, i.Interaction.GuildID)

		responseMessage, err := s.InteractionResponse(i.Interaction)
		if err != nil {
			fmt.Println("Error getting interaction response message:", err)
		}

		threadID := i.Interaction.ID
		responseMessageID := ""
		if responseMessage != nil && responseMessage.ID != "" {
			threadID = responseMessage.ID
			responseMessageID = responseMessage.ID
		}

		err = persistance.SaveConversationMessage(persistance.ConversationMessage{
			ThreadId:    threadID,
			ChannelId:   i.Interaction.ChannelID,
			GuildId:     i.Interaction.GuildID,
			CommandName: "bot-gpt",
			Role:        "user",
			Content:     prompt,
		})
		if err != nil {
			fmt.Println("Error saving conversation user message:", err)
		}

		err = persistance.SaveConversationMessage(persistance.ConversationMessage{
			ThreadId:    threadID,
			MessageId:   responseMessageID,
			ChannelId:   i.Interaction.ChannelID,
			GuildId:     i.Interaction.GuildID,
			CommandName: "bot-gpt",
			Role:        "assistant",
			Content:     assistantResponse,
		})
		if err != nil {
			fmt.Println("Error saving conversation assistant message:", err)
		}

	}
}

func (b *BotGPT) HelpString() string {
	return "The `/bot-gpt` command allows you to prompt OpenAI's GPT-3 or GPT-4 chat model. You can ask it whatever as part of the `prompt` and once it generates a response, it will update the message with what came back. This is slower than the `/bot` command due to the chat model being more complex."
}

func (b *BotGPT) CommandCost() int {
	return 0
}

func ParseGPTSlashCommand(prompt string) (string, string) {
	respTxt := external.GetOpenAIGPTResponse(prompt)
	displayText := "Request: " + prompt + " " + respTxt
	return displayText, respTxt
}
