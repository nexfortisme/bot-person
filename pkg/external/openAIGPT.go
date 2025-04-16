package external

import (
	"encoding/json"
	"fmt"
	"io"
	"main/pkg/logging"
	"main/pkg/util"
	"net/http"
	"strings"
)

func GetOpenAIGPTResponse(prompt string) string {
	client := &http.Client{}

	dataTemplate := `{
		"model": %s,
		"messages": [{"role": "system", "content": "You generate responses no longer than 1750 characters long."}, {"role": "user", "content": "%s"}]
	}`

	data := fmt.Sprintf(dataTemplate, util.EscapeQuotes(util.GetOpenAIModel()), util.EscapeQuotes(prompt))

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

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
