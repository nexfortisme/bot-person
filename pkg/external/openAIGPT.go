package external

import (
	"encoding/json"
	"fmt"
	"io"
	external "main/pkg/external/models"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"main/pkg/util"
	"net/http"
	"strings"
)

var (
	model = "gpt-4"
)

func GetOpenAIGPTResponse(prompt string) string {
	client := &http.Client{}

	dataTemplate := `{
		"model": "%s",
		"messages": [{"role": "system", "content": "You generate responses no longer than 1750 characters long."}, {"role": "user", "content": "%s"}]
	}`
	data := fmt.Sprintf(dataTemplate, model, prompt)

	logging.LogEvent(loggingType.EXTERNAL_GPT_REQUEST, data, "System", "System", nil)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "Error Creating OpenAI Request", "System", "System", nil)
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	buf, _ := io.ReadAll(resp.Body)
	rspOAI := external.OpenAIGPTResponse{}
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "OpenAI Response Empty", "System", "System", nil)
		return "I'm sorry, I don't understand?"
	} else {
		response := rspOAI.Choices[0].Message.Content
		logging.LogEvent(loggingType.EXTERNAL_GPT_RESPONSE, response, "System", "System", nil)
		return response
	}
}
