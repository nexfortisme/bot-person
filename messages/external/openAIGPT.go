package external

import (
	"encoding/json"
	"fmt"
	"io"
	"main/logging"
	"net/http"
	"strings"
)

var (
	model = "gpt-4"
)

func GetOpenAIGPTResponse(prompt string, openAIKey string) string {
	client := &http.Client{}

	dataTemplate := `{
		"model": "%s",
		"messages": [{"role": "user", "content": "%s"}]
	}`

	data := fmt.Sprintf(dataTemplate, model, prompt)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIKey)

	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

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

func SetGPT4() {
	model = "gpt-4"
}

func SetGPT3() {
	model = "gpt-3.5-turbo-0301"
}
