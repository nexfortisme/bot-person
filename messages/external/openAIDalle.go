package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"main/logging"
	"net/http"
	"strings"
)

func GetDalleResponse(prompt string, openAIKey string) (string, error) {
	client := &http.Client{}

	dataTemplate := `{
		"prompt": "%s",
		"n": 1,
		"size": "1024x1024"
	  }`
	data := fmt.Sprintf(dataTemplate, prompt)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/images/generations", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIKey)

	resp, _ := client.Do(req)

	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later.", errors.New("API Error")
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	var rspOAI DalleResponse
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return "", errors.New("Parse Error")
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Data) == 0 {
		return "I'm sorry, I don't understand? (Most likely picked up by OpenAi query filter).", errors.New("API Response Error")
	} else {
		return rspOAI.Data[0].URL, nil
	}
}
