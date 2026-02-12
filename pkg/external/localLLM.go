package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"main/pkg/logging"

	"net/http"
)

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
	client := &http.Client{}
	model := "gemma3-qat"
	requestBody := ""
	responseBody := ""
	statusCode := 0
	var requestErr error
	defer func() {
		logLocalLLMRequest("local_llm", userId, model, requestBody, responseBody, statusCode, requestErr)
	}()

	systemPrompt := "Have your response be funny]. Include a joke at the expense of the user, or be sarcastic. Keep your responses short and to the point."
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
		requestErr = err
		logging.LogError("Error creating request body for local LLM chat completions")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	requestBody = string(body)

	req, err := http.NewRequest(http.MethodPost, localLLMChatCompletionsEndpoint, bytes.NewReader(body))
	if err != nil {
		requestErr = err
		logging.LogError("Error creating POST request")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		requestErr = err
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	if resp == nil {
		requestErr = fmt.Errorf("local LLM response was nil")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	defer resp.Body.Close()
	statusCode = resp.StatusCode

	buf, err := io.ReadAll(resp.Body)
	responseBody = string(buf)
	if err != nil {
		requestErr = err
		return "Error Contacting Local LLM API. Please Try Again Later."
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		requestErr = fmt.Errorf("local LLM returned status code %d", resp.StatusCode)
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	rspOAI := OpenAIGPTResponse{}
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		requestErr = err
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
