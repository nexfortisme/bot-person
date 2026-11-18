package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"main/pkg/logging"
	"main/pkg/util"
	"main/system-prompts"
	"net/http"
	"strings"
	"time"
)

const (
	openRouterChatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"
	openRouterImagesURL          = "https://openrouter.ai/api/v1/images"
	openRouterImageModel         = "openai/gpt-image-2.6-sunburst"
	openRouterErrorMessage       = "Error contacting OpenRouter. Please try again later."
	retortSystemPrompt           = "Have the response be a retort to the user's message. Be funny and sarcastic."
)

var openRouterImageClient = &http.Client{Timeout: 5 * time.Minute}

type openRouterPlugin struct {
	ID string `json:"id"`
}

type openRouterChatRequest struct {
	Model    string             `json:"model"`
	Messages any                `json:"messages"`
	Stream   bool               `json:"stream,omitempty"`
	Plugins  []openRouterPlugin `json:"plugins,omitempty"`
}

type openRouterChatEnvelope struct {
	Error   *Error `json:"error"`
	Choices []struct {
		Message struct {
			Content json.RawMessage `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// NormalizeOpenRouterModel turns a configured model name into an OpenRouter slug.
// Bare OpenAI ids get an openai/ prefix.
func NormalizeOpenRouterModel(model string) string {
	model = strings.TrimSpace(model)
	if model == "" || strings.Contains(model, "/") || strings.HasPrefix(model, "~") {
		return model
	}
	return "openai/" + model
}

// ImageGenerationModel is the OpenRouter image model. GPT Image 2.6 Sunburst is
// the default, including when the env value is still a chat model id.
func ImageGenerationModel() string {
	raw := strings.TrimSpace(util.GetImageGenerationModel())
	switch strings.ToLower(raw) {
	case "", "gpt-5.4", "openai/gpt-5.4":
		return openRouterImageModel
	default:
		return NormalizeOpenRouterModel(raw)
	}
}

func resolveChatCommand(command string) (string, string) {
	canonical := canonicalCommand(command)
	if route, ok := RouteByCommand(canonical); ok {
		return route.Model, systemprompts.ByMention(route.Mention)
	}
	return NormalizeOpenRouterModel(util.GetOpenAIModel()), systemprompts.ByMention("chatgpt")
}

func withSystemPrompt(system string, messages []OpenAIChatMessage) []OpenAIChatMessage {
	if strings.TrimSpace(system) == "" {
		return messages
	}
	out := make([]OpenAIChatMessage, 0, len(messages)+1)
	out = append(out, OpenAIChatMessage{Role: "system", Content: system})
	out = append(out, messages...)
	return out
}

// StreamBotCommand streams a chat completion for a bot command name.
func StreamBotCommand(command string, messages []OpenAIChatMessage, webSearch bool, onDelta func(string)) (string, error) {
	model, system := resolveChatCommand(command)
	return StreamOpenRouterChat(model, withSystemPrompt(system, messages), webSearch, onDelta)
}

// CompleteBotCommand returns a full chat completion for a bot command name.
func CompleteBotCommand(command string, messages []OpenAIChatMessage, webSearch bool) string {
	model, system := resolveChatCommand(command)
	return completeOrErrorMessage(model, withSystemPrompt(system, messages), webSearch)
}

func StreamOpenRouterChat(model string, messages []OpenAIChatMessage, webSearch bool, onDelta func(string)) (string, error) {
	payload := openRouterChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   true,
	}
	if webSearch {
		payload.Plugins = []openRouterPlugin{{ID: "web"}}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	response, _, err := streamChatCompletions(
		openRouterChatCompletionsURL,
		body,
		"Bearer "+util.GetOpenRouterAPIKey(),
		onDelta,
	)
	if err != nil {
		return "", err
	}
	return response, nil
}

func CompleteOpenRouterChat(model string, messages []OpenAIChatMessage, webSearch bool) (string, error) {
	payload := openRouterChatRequest{
		Model:    model,
		Messages: messages,
	}
	if webSearch {
		payload.Plugins = []openRouterPlugin{{ID: "web"}}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, openRouterChatCompletionsURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	setOpenRouterHeaders(req)

	resp, err := llmHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("received nil response from OpenRouter")
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("openrouter status %d: %s", resp.StatusCode, truncateForLog(string(buf)))
	}

	text, err := assistantTextFromBody(buf)
	if err != nil {
		return "", err
	}
	return text, nil
}

func completeOrErrorMessage(model string, messages []OpenAIChatMessage, webSearch bool) string {
	text, err := CompleteOpenRouterChat(model, messages, webSearch)
	if err != nil {
		logging.LogError(fmt.Sprintf("OpenRouter request failed for %s: %v", model, err))
		return openRouterErrorMessage
	}
	if strings.TrimSpace(text) == "" {
		return "I'm sorry, I don't understand?"
	}
	return text
}

func setOpenRouterHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+util.GetOpenRouterAPIKey())
	req.Header.Set("X-OpenRouter-Title", "Bot Person")
}

func assistantTextFromBody(buf []byte) (string, error) {
	envelope := openRouterChatEnvelope{}
	if err := json.Unmarshal(buf, &envelope); err != nil {
		return "", err
	}
	if envelope.Error != nil && strings.TrimSpace(envelope.Error.Message) != "" {
		return "", fmt.Errorf("%s", envelope.Error.Message)
	}
	if len(envelope.Choices) == 0 {
		return "", nil
	}
	return contentToString(envelope.Choices[0].Message.Content), nil
}

func contentToString(raw json.RawMessage) string {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 || string(trimmed) == "null" {
		return ""
	}

	var text string
	if err := json.Unmarshal(trimmed, &text); err == nil {
		return text
	}

	var parts []struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(trimmed, &parts); err == nil {
		builder := strings.Builder{}
		for _, part := range parts {
			builder.WriteString(part.Text)
		}
		return builder.String()
	}

	return ""
}

func truncateForLog(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 500 {
		return value[:500]
	}
	return value
}
