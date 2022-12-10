package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetDalleResponse(prompt string, openAIKey string) (string, error) {
	httpClient := &http.Client{}

	requestDataTemplate := `{
		"prompt": "%s",
		"n": 1,
		"size": "1024x1024"
	  }`
	requestData := fmt.Sprintf(requestDataTemplate, prompt)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/images/generations", strings.NewReader(requestData))
	if err != nil {
		return "Error creating POST request.", errors.New("POST Request Error")
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+openAIKey)

	httpResponse, _ := httpClient.Do(postRequest)

	if httpResponse == nil {
		return "Error Contacting OpenAI API. Please Try Again Later.", errors.New("API Error")
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)
	var openAIResponse DalleResponse
	err = json.Unmarshal([]byte(string(responseBuffer)), &openAIResponse)
	if err != nil {
		return "Error Parsing Response", errors.New("paree response error")
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(openAIResponse.Data) == 0 {
		// fmt.Println(responseBuffer)
		// fmt.Println(openAIResponse)
		return "I'm sorry, I don't understand? (Most likely picked up by OpenAi query filter).", errors.New("API Response Error")
	} else {
		return openAIResponse.Data[0].URL, nil
	}
}
