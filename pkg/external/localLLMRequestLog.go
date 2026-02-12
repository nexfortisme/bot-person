package external

import (
	"fmt"
	"main/pkg/logging"
	"main/pkg/persistance"
)

const localLLMChatCompletionsEndpoint = "http://localhost:12434/engines/v1/chat/completions"

func logLocalLLMRequest(
	requestType string,
	userID string,
	model string,
	requestBody string,
	responseBody string,
	statusCode int,
	requestErr error,
) {
	errorMessage := ""
	if requestErr != nil {
		errorMessage = requestErr.Error()
	}

	err := persistance.SaveLocalLLMLog(persistance.LocalLLMLog{
		RequestType:  requestType,
		UserId:       userID,
		Model:        model,
		Endpoint:     localLLMChatCompletionsEndpoint,
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	})
	if err != nil {
		logging.LogError(fmt.Sprintf("Error saving local LLM request log: %v", err))
	}
}
