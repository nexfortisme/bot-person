package external

func StreamOpenAIGPTResponse(prompt string, onDelta func(string)) (string, error) {
	return StreamOpenAIGPTResponseWithChatMessages([]OpenAIChatMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}, onDelta)
}

func StreamOpenAIGPTResponseWithChatMessages(messages []OpenAIChatMessage, onDelta func(string)) (string, error) {
	return StreamBotCommand("bot-gpt", messages, false, onDelta)
}

func GetOpenAIGPTResponse(prompt string) string {
	return GetOpenAIGPTResponseWithMessages([]OpenAIGPTMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	})
}

func GetOpenAIGPTResponseWithMessages(messages []OpenAIGPTMessage) string {
	return GetOpenAIGPTResponseWithChatMessages(gptMessagesToChat(messages))
}

func GetOpenAIGPTResponseWithChatMessages(messages []OpenAIChatMessage) string {
	return CompleteBotCommand("bot-gpt", messages, false)
}

func StreamOpenAIGPTResponseWithWebSearch(prompt string, onDelta func(string)) (string, error) {
	return StreamOpenAIGPTResponseWithChatMessagesAndWebSearch([]OpenAIChatMessage{
		{Role: "user", Content: prompt},
	}, onDelta)
}

func StreamOpenAIGPTResponseWithChatMessagesAndWebSearch(messages []OpenAIChatMessage, onDelta func(string)) (string, error) {
	return StreamBotCommand("bot-gpt", messages, true, onDelta)
}

func GetOpenAIGPTResponseWithWebSearch(prompt string) string {
	return GetOpenAIGPTResponseWithChatMessagesAndWebSearch([]OpenAIChatMessage{
		{Role: "user", Content: prompt},
	})
}

func GetOpenAIGPTResponseWithChatMessagesAndWebSearch(messages []OpenAIChatMessage) string {
	return CompleteBotCommand("bot-gpt", messages, true)
}

func gptMessagesToChat(messages []OpenAIGPTMessage) []OpenAIChatMessage {
	requestMessages := make([]OpenAIChatMessage, 0, len(messages))
	for _, message := range messages {
		requestMessages = append(requestMessages, OpenAIChatMessage{
			Role:    message.Role,
			Content: message.Content,
		})
	}
	return requestMessages
}
