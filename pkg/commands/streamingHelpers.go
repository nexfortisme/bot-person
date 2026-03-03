package commands

import (
	"fmt"
	"main/pkg/util"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	interactionResponseLookupAttempts = 15
	interactionResponseLookupDelay    = 200 * time.Millisecond
	streamInProgressSuffix            = "\n\n[Response still generating...]"
	longResponseFileSuffix            = "\n\n[Full response attached in too-long.txt]"
)

func getInteractionResponseMessageWithRetry(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.Message {
	for attempt := 0; attempt < interactionResponseLookupAttempts; attempt++ {
		responseMessage, err := s.InteractionResponse(i.Interaction)
		if err == nil && responseMessage != nil && responseMessage.ID != "" {
			return responseMessage
		}

		time.Sleep(interactionResponseLookupDelay)
	}

	return nil
}

func formatPromptResponse(prompt string, assistantResponse string) string {
	trimmedPrompt := strings.TrimSpace(prompt)
	trimmedAssistantResponse := strings.TrimSpace(assistantResponse)

	if trimmedPrompt == "" {
		if trimmedAssistantResponse == "" {
			return "Thinking..."
		}
		return trimmedAssistantResponse
	}

	if trimmedAssistantResponse == "" {
		return "Request: " + trimmedPrompt + "\n\nThinking..."
	}

	return "Request: " + trimmedPrompt + "\n\n" + trimmedAssistantResponse
}

func formatPromptResponseInProgress(prompt string, assistantResponse string) string {
	return appendSuffixWithinDiscordLimit(formatPromptResponse(prompt, assistantResponse), streamInProgressSuffix)
}

func formatLongResponsePreview(content string) string {
	return appendSuffixWithinDiscordLimit(content, longResponseFileSuffix)
}

func appendSuffixWithinDiscordLimit(content string, suffix string) string {
	if len(suffix) >= 2000 {
		return util.TruncateForDiscord(suffix)
	}

	baseBudget := 2000 - len(suffix)
	baseContent := content
	if len(baseContent) > baseBudget {
		if baseBudget > 3 {
			baseContent = baseContent[:baseBudget-3] + "..."
		} else {
			baseContent = baseContent[:baseBudget]
		}
	}

	return baseContent + suffix
}

func streamInteractionResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	initialContent string,
	formatDisplay func(assistantResponse string) string,
	streamFn func(onDelta func(string)) (string, error),
) (string, error) {
	return util.StreamResponseWithThrottledEdits(
		initialContent,
		formatDisplay,
		streamFn,
		func(content string) error {
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Content: &content,
			})
			if err != nil {
				fmt.Println("Error editing interaction response while streaming:", err)
			}
			return err
		},
	)
}
