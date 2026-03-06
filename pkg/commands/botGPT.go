package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"strings"

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
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "web-search",
				Description: "Enable web search to get up-to-date information.",
				Required:    false,
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

	// Pulling the propt out of the optionsMap
	if option, ok := optionMap["prompt"]; ok {

		prompt := option.StringValue()
		logging.LogEvent(eventType.COMMAND_BOT_GPT, i.Interaction.Member.User.ID, prompt, i.Interaction.GuildID)

		useWebSearch := false
		if webSearchOption, ok := optionMap["web-search"]; ok {
			useWebSearch = webSearchOption.BoolValue()
		}

		initialResponse := formatPromptResponseInProgress(prompt, "")
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: initialResponse,
			},
		})
		if err != nil {
			fmt.Println("Error creating interaction response:", err)
			return
		}

		responseMessage := getInteractionResponseMessageWithRetry(s, i)
		threadID := i.Interaction.ID
		responseMessageID := ""
		if responseMessage != nil && responseMessage.ID != "" {
			threadID = responseMessage.ID
			responseMessageID = responseMessage.ID
		}

		if !util.TryStartThreadResponse(threadID) {
			busyResponse := "Please wait for the current response in this thread to finish before sending another reply."
			_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &busyResponse,
			})
			return
		}
		defer util.FinishThreadResponse(threadID)

		formatFn := func(assistant string) string {
			return formatPromptResponseInProgress(prompt, assistant)
		}
		streamFn := func(onDelta func(string)) (string, error) {
			if useWebSearch {
				return external.StreamOpenAIGPTResponseWithWebSearch(prompt, onDelta)
			}
			return external.StreamOpenAIGPTResponse(prompt, onDelta)
		}

		var assistantResponse string
		var streamErr error
		if useWebSearch {
			assistantResponse, streamErr = streamInteractionResponseSuppressEmbeds(s, i, initialResponse, formatFn, streamFn)
		} else {
			assistantResponse, streamErr = streamInteractionResponse(s, i, initialResponse, formatFn, streamFn)
		}

		if streamErr != nil || strings.TrimSpace(assistantResponse) == "" {
			if streamErr != nil {
				fmt.Println("Error streaming /bot-gpt response:", streamErr)
			}

			if useWebSearch {
				assistantResponse = external.GetOpenAIGPTResponseWithWebSearch(prompt)
			} else {
				assistantResponse = external.GetOpenAIGPTResponse(prompt)
			}
			if strings.TrimSpace(assistantResponse) == "" {
				assistantResponse = "I'm sorry, I don't understand?"
			}

			fallbackContent := util.TruncateForDiscord(formatPromptResponse(prompt, assistantResponse))
			var fallbackErr error
			if useWebSearch {
				// Clear SUPPRESS_EMBEDS flag (flags: 0) so sources render on the final message.
				fallbackErr = interactionResponseEditWithFlags(s, i, fallbackContent, 0)
			} else {
				_, fallbackErr = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &fallbackContent,
				})
			}
			if fallbackErr != nil {
				fmt.Println("Error editing interaction response:", fallbackErr)
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong.",
				})
				return
			}
		}

		displayResponse := formatPromptResponse(prompt, assistantResponse)
		logging.LogEvent(eventType.EXTERNAL_GPT_RESPONSE, i.Interaction.Member.User.ID, assistantResponse, i.Interaction.GuildID)

		if len(displayResponse) > 2000 {
			fileObj := util.HandleTooLongResponseWithFileName(displayResponse, "too-long.txt")
			fileMessage := formatLongResponsePreview(displayResponse)

			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &fileMessage,
				Files:   []*discordgo.File{fileObj},
			})

			if err != nil {
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong while sending the full response as a file.",
				})
				return
			}
		} else if useWebSearch {
			// Clear SUPPRESS_EMBEDS flag (flags: 0) so sources render on the final message.
			err := interactionResponseEditWithFlags(s, i, displayResponse, 0)
			if err != nil {
				fmt.Println("Error editing interaction response:", err)
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong.",
				})
				return
			}
		} else {
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &displayResponse,
			})
			if err != nil {
				fmt.Println("Error editing interaction response:", err)
				s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong.",
				})
				return
			}
		}

		if responseMessageID == "" {
			responseMessage = getInteractionResponseMessageWithRetry(s, i)
			if responseMessage != nil && responseMessage.ID != "" {
				threadID = responseMessage.ID
				responseMessageID = responseMessage.ID
			}
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
	return "The `/bot-gpt` command streams GPT output as it is generated. Use the `web-search` option to enable real-time web search for up-to-date information. Reply to a `/bot-gpt` message to continue the same conversation thread. While one response is in progress, additional thread replies are temporarily blocked."
}

func (b *BotGPT) CommandCost() int {
	return 0
}
