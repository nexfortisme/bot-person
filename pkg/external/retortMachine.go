package external

import (
	"bytes"
	"encoding/json"
	"io"

	"main/pkg/logging"

	"net/http"
)

func GetRetortMachineResponse(prompt string, userId string) string {
	return GetRetortMachineResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, userId)
}

func GetRetortMachineResponseWithMessages(messages []OpenAIGPTMessage, _ string) string {
	client := &http.Client{}

	systemPrompt := "Have the response be a retort to the user's message. Be funny and sarcastic."
	requestMessages := make([]OpenAIGPTMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIGPTMessage{
		Role:    "system",
		Content: systemPrompt,
	})
	requestMessages = append(requestMessages, messages...)

	payload := chatCompletionsRequest{
		Model:    "smollm2",
		Messages: requestMessages,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		logging.LogError("Error creating request body for local LLM chat completions")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:12434/engines/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		logging.LogError("Error creating POST request")
		return "Error Contacting Local LLM API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting Local LLM API. Please Try Again Later."
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
