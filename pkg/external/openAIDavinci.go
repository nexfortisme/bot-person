package external

import (
	"main/pkg/persistance"
	attribute "main/pkg/persistance/eums"
	"main/pkg/util"
)

func GetOpenAIResponse(prompt string, userId string) string {
	return GetOpenAIResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, userId)
}

func GetOpenAIResponseWithMessages(messages []OpenAIGPTMessage, userId string) string {
	prePrompt, err := persistance.GetUserAttribute(userId, attribute.BOT_PREPROMPT)
	if err != nil {
		prePrompt = "You are a whimsical and dear friend. You respond to any inquiries with a level of spontaneity and randomness. You don't take anything too seriously and are not afraid to 'shoot from the hip' so to speak when responding to people."
	}

	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: prePrompt,
	})
	requestMessages = append(requestMessages, gptMessagesToChat(messages)...)

	return completeOrErrorMessage(NormalizeOpenRouterModel(util.GetBotOpenAIModel()), requestMessages, false)
}
