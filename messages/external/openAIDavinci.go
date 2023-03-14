package external

import (
	"encoding/json"
	"fmt"
	"io"
	"main/logging"
	"net/http"
	"strings"
)

func GetOpenAIResponse(prompt string, openAIKey string) string {
	client := &http.Client{}

	// dataTemplate := `{
	// 	"model": "gpt-3.5-turbo",
	// 	"prompt": "%s",
	// 	"temperature": 0.7,
	// 	"max_tokens": 256,
	// 	"top_p": 1,
	// 	"frequency_penalty": 0,
	// 	"presence_penalty": 0
	//   }`
	dataTemplate := `
		"model": "gpt-3.5-turbo",
		"messages": [{"role": "user", "content": "%s"}]
	`
	data := fmt.Sprintf(dataTemplate, prompt)

	fmt.Println(data)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIKey)

	fmt.Println(req.Body)

	resp, _ := client.Do(req)

	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	buf, _ := io.ReadAll(resp.Body)
	var rspOAI OpenAIResponse
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	fmt.Println(string(buf))

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Message.Content
	}
}
