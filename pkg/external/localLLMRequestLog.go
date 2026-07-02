package external

import (
	"fmt"
	"main/pkg/logging"
	"main/pkg/persistance"
	"main/pkg/util"
)



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
		Endpoint:     util.GetLocalLLMChatCompletionsEndpoint(),
		RequestBody:  requestBody,
		ResponseBody: responseBody,
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	})
	if err != nil {
		logging.LogError(fmt.Sprintf("Error saving local LLM request log: %v", err))
	}
}
