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

func GetLocalLLMResponse(prompt string, userId string) string {
	client := &http.Client{}

	// prePrompt, err := persistance.GetUserAttribute(userId, attribute.BOT_PREPROMPT)
	// if err != nil {
	// 	prePrompt = "You are a whimsical and dear friend. You respond to any inquiries with a level of spontaneity and randomness. You don't take anything too seriously and are not afraid to 'shoot from the hip' so to speak when responding to people."
	// }

	systemPrompt := "Have your response be funny, even if it is not relevant. Include a joke at the expense of the user."

	dataTemplate := `{
		"model": "gemma3-qat",
		"messages": [{"role": "system", "content": "%s"}, {"role": "user", "content": "%s"}]
	}`

	data := fmt.Sprintf(dataTemplate, util.EscapeQuotes(systemPrompt), util.EscapeQuotes(prompt))

	fmt.Println(data)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:12434/engines/v1/chat/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	resp, _ := client.Do(req)
	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	fmt.Println(resp)

	buf, _ := io.ReadAll(resp.Body)
	rspOAI := OpenAIGPTResponse{}
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	fmt.Println(rspOAI)

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Message.Content
	}
}
