package external

import (
	"bufio"
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"main/pkg/logging"
	"main/pkg/util"
	"net/http"
	"strings"
)

//go:embed rocky_system_prompt.md
var defaultGPTSystemPrompt string

type chatCompletionsRequest struct {
	Model            string `json:"model"`
	Messages         any    `json:"messages"`
	Reasoning_Effort string `json:"reasoning_effort,omitempty"`
}

type openAIWebSearchTool struct {
	Type string `json:"type"`
}

type responsesAPIRequest struct {
	Model        string                `json:"model"`
	Input        any                   `json:"input"`
	Instructions string                `json:"instructions,omitempty"`
	Tools        []openAIWebSearchTool `json:"tools"`
}

type responsesAPIStreamRequest struct {
	Model        string                `json:"model"`
	Input        any                   `json:"input"`
	Instructions string                `json:"instructions,omitempty"`
	Tools        []openAIWebSearchTool `json:"tools"`
	Stream       bool                  `json:"stream"`
}

type responsesAPIOutputTextDelta struct {
	Type  string `json:"type"`
	Delta string `json:"delta"`
}

func openAIWebSearchTools() []openAIWebSearchTool {
	return []openAIWebSearchTool{{Type: "web_search"}}
}

func StreamOpenAIGPTResponse(prompt string, onDelta func(string)) (string, error) {
	return StreamOpenAIGPTResponseWithChatMessages([]OpenAIChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, onDelta)
}

func StreamOpenAIGPTResponseWithChatMessages(messages []OpenAIChatMessage, onDelta func(string)) (string, error) {
	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: defaultGPTSystemPrompt,
	})
	requestMessages = append(requestMessages, messages...)

	payload := streamChatCompletionsRequest{
		Model:    util.GetOpenAIModel(),
		Messages: requestMessages,
		Stream:   true,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	response, _, err := streamChatCompletions(
		"https://api.openai.com/v1/chat/completions",
		requestBody,
		"Bearer "+util.GetOpenAIKey(),
		onDelta,
	)
	if err != nil {
		return "", err
	}

	return response, nil
}

func GetOpenAIGPTResponse(prompt string) string {
	return GetOpenAIGPTResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	})
}

func GetOpenAIGPTResponseWithMessages(messages []OpenAIGPTMessage) string {
	requestMessages := make([]OpenAIChatMessage, 0, len(messages))
	for _, message := range messages {
		requestMessages = append(requestMessages, OpenAIChatMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	return GetOpenAIGPTResponseWithChatMessages(requestMessages)
}

func GetOpenAIGPTResponseWithChatMessages(messages []OpenAIChatMessage) string {
	client := &http.Client{}
	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: defaultGPTSystemPrompt,
	})
	requestMessages = append(requestMessages, messages...)

	payload := chatCompletionsRequest{
		Model:    util.GetOpenAIModel(),
		Messages: requestMessages,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.LogError("Error creating request body for OpenAI chat completions")
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		logging.LogError("Error creating POST request")
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	rspOAI := OpenAIGPTResponse{}
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Message.Content
	}
}

func StreamOpenAIGPTResponseWithWebSearch(prompt string, onDelta func(string)) (string, error) {
	return StreamOpenAIGPTResponseWithChatMessagesAndWebSearch([]OpenAIChatMessage{
		{Role: "user", Content: prompt},
	}, onDelta)
}

func StreamOpenAIGPTResponseWithChatMessagesAndWebSearch(messages []OpenAIChatMessage, onDelta func(string)) (string, error) {
	payload := responsesAPIStreamRequest{
		Model:        util.GetOpenAIModel(),
		Input:        messages,
		Instructions: defaultGPTSystemPrompt,
		Tools:        openAIWebSearchTools(),
		Stream:       true,
	}

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("received nil response from Responses API")
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("responses API request failed with status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return parseResponsesAPIStream(resp.Body, onDelta)
}

func parseResponsesAPIStream(body io.Reader, onDelta func(string)) (string, error) {
	scanner := bufio.NewScanner(body)
	fullResponse := strings.Builder{}

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")
		if payload == "[DONE]" {
			break
		}

		var delta responsesAPIOutputTextDelta
		if err := json.Unmarshal([]byte(payload), &delta); err != nil {
			continue
		}

		if delta.Type == "response.output_text.delta" && delta.Delta != "" {
			fullResponse.WriteString(delta.Delta)
			if onDelta != nil {
				onDelta(delta.Delta)
			}
		}
	}

	return fullResponse.String(), scanner.Err()
}

func GetOpenAIGPTResponseWithWebSearch(prompt string) string {
	return GetOpenAIGPTResponseWithChatMessagesAndWebSearch([]OpenAIChatMessage{
		{Role: "user", Content: prompt},
	})
}

func GetOpenAIGPTResponseWithChatMessagesAndWebSearch(messages []OpenAIChatMessage) string {
	payload := responsesAPIRequest{
		Model:        util.GetOpenAIModel(),
		Input:        messages,
		Instructions: defaultGPTSystemPrompt,
		Tools:        openAIWebSearchTools(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.LogError("Error creating request body for Responses API web search")
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		logging.LogError("Error creating POST request")
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}
	defer resp.Body.Close()

	buf, _ := io.ReadAll(resp.Body)
	rspOAI := Response{}
	err = json.Unmarshal(buf, &rspOAI)
	if err != nil {
		return ""
	}

	for _, output := range rspOAI.Output {
		if output.Type != "message" {
			continue
		}
		for _, content := range output.Content {
			if content.Type == "output_text" && content.Text != "" {
				return content.Text
			}
		}
	}

	return "I'm sorry, I don't understand?"
}
