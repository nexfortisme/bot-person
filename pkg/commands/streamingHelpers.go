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

// webhookEditWithFlags is a raw edit payload that includes Flags without omitempty,
// allowing us to explicitly send flags: 0 to clear SUPPRESS_EMBEDS.
// discordgo.WebhookEdit does not have a Flags field, and discordgo.MessageEdit
// uses omitempty on Flags which prevents sending 0 to clear a flag.
type webhookEditWithFlags struct {
	Content *string `json:"content,omitempty"`
	Flags   int     `json:"flags"`
}

// interactionResponseEditWithFlags patches the original interaction response via
// the webhook endpoint with an explicit flags value.
func interactionResponseEditWithFlags(s *discordgo.Session, i *discordgo.InteractionCreate, content string, flags int) error {
	uri := discordgo.EndpointWebhookMessage(i.Interaction.AppID, i.Interaction.Token, "@original")
	_, err := s.RequestWithBucketID("PATCH", uri,
		&webhookEditWithFlags{Content: &content, Flags: flags},
		discordgo.EndpointWebhookToken("", ""))
	return err
}

// streamInteractionResponseSuppressEmbeds is like streamInteractionResponse but
// suppresses Discord link previews on in-progress edits so they don't flash during
// streaming. Call interactionResponseEditWithFlags with flags=0 after streaming to
// clear the suppression and allow embeds on the final message.
func streamInteractionResponseSuppressEmbeds(
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
			err := interactionResponseEditWithFlags(s, i, content, int(discordgo.MessageFlagsSuppressEmbeds))
			if err != nil {
				fmt.Println("Error editing interaction response while streaming:", err)
			}
			return err
		},
	)
}
