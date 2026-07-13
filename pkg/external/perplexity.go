package external

import (
	"bytes"
	"encoding/json"
	"io"
	"main/pkg/logging"
	"main/pkg/util"
	"net/http"
	"time"
)

func GetPerplexityResponse(originalPrompt string, userPrompt string) PerplexityResponse {

	client := &http.Client{Timeout: 30 * time.Second}
	payload := struct {
		Model    string             `json:"model"`
		Messages []OpenAIGPTMessage `json:"messages"`
	}{
		Model: "sonar",
		Messages: []OpenAIGPTMessage{
			{Role: "system", Content: "Please keep responses concise and to the point. Do not include any additional information or commentary."},
			{Role: "user", Content: originalPrompt},
			{Role: "assistant", Content: "The first message is context, the second message is the user's message pertaining to the first."},
			{Role: "user", Content: userPrompt},
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		logging.LogError("Error creating Perplexity request body")
		return PerplexityResponse{}
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.perplexity.ai/chat/completions", bytes.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
		return PerplexityResponse{}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+util.GetPerplexityAPIKey())

	resp, err := client.Do(req)
	if err != nil || resp == nil {
		return PerplexityResponse{}
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return PerplexityResponse{}
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return PerplexityResponse{}
	}

	respPex := PerplexityResponse{}
	// need to add proper error handling
	err = json.Unmarshal(buf, &respPex)
	if err != nil {
		return PerplexityResponse{}
	}

	return respPex
}
