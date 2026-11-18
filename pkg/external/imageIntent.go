package external

import (
	"strings"

	"main/pkg/util"
)

const imageEditIntentSystemPrompt = "You are an intent classifier for a Discord bot. " +
	"The user sent a message about an image. Decide if they want the image edited, transformed, or a new version generated " +
	"(examples: \"make it blue\", \"add a hat\", \"turn this into a cartoon\", \"remove the background\") " +
	"or if they are just asking about or commenting on the image " +
	"(examples: \"what is this?\", \"describe this\", \"lol nice\"). " +
	"Respond with exactly one word: EDIT or CHAT."

// DetectImageEditIntent asks OpenRouter whether the user's message is a
// request to edit an image. Any failure falls back to false so a broken
// classifier never blocks a normal chat response.
func DetectImageEditIntent(prompt string) bool {
	text, err := CompleteOpenRouterChat(
		NormalizeOpenRouterModel(util.GetRetortMachineModel()),
		[]OpenAIChatMessage{
			{Role: "system", Content: imageEditIntentSystemPrompt},
			{Role: "user", Content: prompt},
		},
		false,
	)
	if err != nil || strings.TrimSpace(text) == "" {
		return false
	}

	verdict := strings.ToUpper(strings.TrimSpace(text))
	return strings.HasPrefix(verdict, "EDIT")
}
