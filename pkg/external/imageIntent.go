package external

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"main/pkg/util"
)

const imageEditIntentSystemPrompt = "You are an intent classifier for a Discord bot. " +
	"The user sent a message about an image. Decide if they want the image edited, transformed, or a new version generated " +
	"(examples: \"make it blue\", \"add a hat\", \"turn this into a cartoon\", \"remove the background\") " +
	"or if they are just asking about or commenting on the image " +
	"(examples: \"what is this?\", \"describe this\", \"lol nice\"). " +
	"Respond with exactly one word: EDIT or CHAT."

// DetectImageEditIntent asks the local LLM whether the user's message is a
// request to edit an image. Any failure falls back to false so a broken
// classifier never blocks a normal chat response.
func DetectImageEditIntent(prompt string) bool {
	payload := chatCompletionsRequest{
		Model: util.GetRetortMachineModel(),
		Messages: []OpenAIChatMessage{
			{Role: "system", Content: imageEditIntentSystemPrompt},
			{Role: "user", Content: prompt},
		},
		Reasoning_Effort: "none",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	req, err := http.NewRequest(http.MethodPost, util.GetLocalLLMChatCompletionsEndpoint(), bytes.NewReader(body))
	if err != nil {
		return false
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := llmHTTPClient.Do(req)
	if err != nil || resp == nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return false
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	rspOAI := OpenAIGPTResponse{}
	if err := json.Unmarshal(buf, &rspOAI); err != nil || len(rspOAI.Choices) == 0 {
		return false
	}

	verdict := strings.ToUpper(strings.TrimSpace(rspOAI.Choices[0].Message.Content))
	return strings.HasPrefix(verdict, "EDIT")
}
