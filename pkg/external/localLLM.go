package external

import (
	"bytes"
	"encoding/json"
	"io"

	"main/pkg/logging"
	"main/pkg/util"

	"net/http"
)

func StreamLocalLLMResponse(prompt string, userId string, onDelta func(string)) (string, error) {
	return StreamLocalLLMResponseWithChatMessages([]OpenAIChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, userId, onDelta)
}

func StreamLocalLLMResponseWithChatMessages(messages []OpenAIChatMessage, userId string, onDelta func(string)) (string, error) {
	model := util.GetLocalLLMModel()

	systemPrompt := "Have your response be funny. Include a joke at the expense of the user, or be sarcastic. Keep your responses short and to the point."
	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})
	requestMessages = append(requestMessages, messages...)

	payload := streamChatCompletionsRequest{
		Model:            model,
		Messages:         requestMessages,
		Stream:           true,
		Reasoning_Effort: "none",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	assistantResponse, _, err := streamChatCompletions(
		util.GetLocalLLMChatCompletionsEndpoint(),
		body,
		"",
		onDelta,
	)
	if err != nil {
		return "", err
	}

	return assistantResponse, nil
}

func GetLocalLLMResponse(prompt string, userId string) string {
	return GetLocalLLMResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, userId)
}

func GetLocalLLMResponseWithMessages(messages []OpenAIGPTMessage, userId string) string {
	requestMessages := make([]OpenAIChatMessage, 0, len(messages))
	for _, message := range messages {
		requestMessages = append(requestMessages, OpenAIChatMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}

	return GetLocalLLMResponseWithChatMessages(requestMessages, userId)
}

func GetLocalLLMResponseWithChatMessages(messages []OpenAIChatMessage, userId string) string {
	client := llmHTTPClient
	model := util.GetLocalLLMModel()

	systemPrompt := "Have your response be funny. Include a joke at the expense of the user, or be sarcastic. Keep your responses short and to the point."
	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: systemPrompt,
	})
	requestMessages = append(requestMessages, messages...)

	payload := chatCompletionsRequest{
		Model:    model,
		Messages: requestMessages,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.LogError("Error creating request body for local LLM chat completions")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	req, err := http.NewRequest(http.MethodPost, util.GetLocalLLMChatCompletionsEndpoint(), bytes.NewReader(body))
	if err != nil {
		logging.LogError("Error creating POST request")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	if resp == nil {
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

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
