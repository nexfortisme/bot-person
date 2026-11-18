package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// RoutedBot is a slash command that chats through one OpenRouter model.
type RoutedBot struct {
	command string
}

func NewRoutedBot(command string) *RoutedBot {
	if _, ok := external.RouteByCommand(command); !ok {
		panic("unknown routed bot " + command)
	}
	return &RoutedBot{command: command}
}

func RoutedBots() []*RoutedBot {
	return []*RoutedBot{
		NewRoutedBot("bot-gemini"),
		NewRoutedBot("bot-muse"),
		NewRoutedBot("bot-claude"),
		NewRoutedBot("bot-grok"),
		NewRoutedBot("bot-bob"),
	}
}

func (b *RoutedBot) CommandName() string {
	return b.command
}

func (b *RoutedBot) ApplicationCommand() *discordgo.ApplicationCommand {
	route, _ := external.RouteByCommand(b.command)
	options := []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "prompt",
			Description: "The prompt to send.",
			Required:    true,
		},
	}
	if route.NativeWebSearch {
		options = append(options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "web-search",
			Description: "Enable web search to get up-to-date information.",
			Required:    false,
		})
	}
	return &discordgo.ApplicationCommand{
		Name:        route.Command,
		Description: route.Description,
		Options:     options,
	}
}

func (b *RoutedBot) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	executeRoutedBot(s, i, b.command)
}

func (b *RoutedBot) HelpString() string {
	route, ok := external.RouteByCommand(b.command)
	if !ok {
		return "Command not found."
	}
	return route.Help
}

func (b *RoutedBot) CommandCost() int {
	return 0
}

func executeRoutedBot(s *discordgo.Session, i *discordgo.InteractionCreate, commandName string) {
	route, ok := external.RouteByCommand(commandName)
	if !ok {
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	option, ok := optionMap["prompt"]
	if !ok {
		return
	}
	prompt := strings.TrimSpace(option.StringValue())
	if prompt == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Give me a prompt to answer.",
			},
		})
		return
	}

	useWebSearch := false
	if route.NativeWebSearch {
		if webSearchOption, ok := optionMap["web-search"]; ok {
			useWebSearch = webSearchOption.BoolValue()
		}
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

	userMessage := []external.OpenAIChatMessage{{
		Role:    "user",
		Content: prompt,
	}}
	formatFn := func(assistant string) string {
		return formatPromptResponseInProgress(prompt, assistant)
	}
	streamFn := func(onDelta func(string)) (string, error) {
		return external.StreamBotCommand(commandName, userMessage, useWebSearch, onDelta)
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
			fmt.Printf("Error streaming /%s response: %v\n", commandName, streamErr)
		}

		assistantResponse = external.CompleteBotCommand(commandName, userMessage, useWebSearch)
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

	if len(displayResponse) > 2000 {
		fileObj := util.HandleTooLongResponseWithFileName(displayResponse, "too-long.txt")
		fileMessage := formatLongResponsePreview(displayResponse)

		_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &fileMessage,
			Files:   []*discordgo.File{fileObj},
		})
		if err != nil {
			fmt.Println("Error editing interaction response:", err)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
		}
	} else if useWebSearch {
		// Clear SUPPRESS_EMBEDS flag (flags: 0) so sources render on the final message.
		err := interactionResponseEditWithFlags(s, i, displayResponse, 0)
		if err != nil {
			fmt.Println("Error editing interaction response:", err)
			s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "Something went wrong.",
			})
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
		CommandName: commandName,
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
		CommandName: commandName,
		Role:        "assistant",
		Content:     assistantResponse,
	})
	if err != nil {
		fmt.Println("Error saving conversation assistant message:", err)
	}
}
