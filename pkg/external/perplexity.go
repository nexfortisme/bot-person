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

func GetPerplexityResponse(originalPrompt string, userPrompt string) PerplexityResponse {

	client := &http.Client{}

	dataTemplate := `{
		"model": "sonar",
		"messages": [
			{"role": "system", "content": "Please keep responses concise and to the point. Do not include any additional information or commentary."},
			{"role": "user", "content": "%s"},
			{"role": "assistant", "content": "The first message is context, the second message is the user's message pertaining to the first."},
			{"role": "user", "content": "%s"}
		]
	}`

	data := fmt.Sprintf(dataTemplate, util.EscapeQuotes(originalPrompt), util.EscapeQuotes(userPrompt))

	req, err := http.NewRequest(http.MethodPost, "https://api.perplexity.ai/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetPerplexityAPIKey())

	resp, _ := client.Do(req)
	if resp == nil {
		return PerplexityResponse{}
	}

	buf, _ := io.ReadAll(resp.Body)
	// fmt.Println(string(buf))

	respPex := PerplexityResponse{}
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &respPex)
	if err != nil {
		return PerplexityResponse{}
	}

	return respPex
}
