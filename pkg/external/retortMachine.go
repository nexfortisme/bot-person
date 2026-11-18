package external

import "main/pkg/util"

func GetRetortMachineResponse(prompt string, userId string) string {
	return GetRetortMachineResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, userId)
}

func GetRetortMachineResponseWithMessages(messages []OpenAIGPTMessage, userId string) string {
	requestMessages := make([]OpenAIChatMessage, 0, len(messages)+1)
	requestMessages = append(requestMessages, OpenAIChatMessage{
		Role:    "system",
		Content: retortSystemPrompt,
	})
	requestMessages = append(requestMessages, gptMessagesToChat(messages)...)

	return completeOrErrorMessage(NormalizeOpenRouterModel(util.GetRetortMachineModel()), requestMessages, false)
}
