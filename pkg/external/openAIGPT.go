package external

import (
	"bytes"
	"encoding/json"
	"io"
	"main/pkg/logging"
	"main/pkg/util"
	"net/http"
)

const defaultGPTSystemPrompt = "You are a single source of truth. Give responses that answer the question asked but don't ask follow up questions."

type chatCompletionsRequest struct {
	Model    string `json:"model"`
	Messages any    `json:"messages"`
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
