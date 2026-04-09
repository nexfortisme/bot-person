package external

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type streamChatCompletionsRequest struct {
	Model            string `json:"model"`
	Messages         any    `json:"messages"`
	Stream           bool   `json:"stream"`
	Reasoning_Effort string `json:"reasoning_effort"`
}

type streamChatCompletionsChunk struct {
	Error   *Error                                `json:"error"`
	Choices []streamChatCompletionsChunkCandidate `json:"choices"`
}

type streamChatCompletionsChunkCandidate struct {
	Text    string `json:"text"`
	Delta   any    `json:"delta"`
	Message any    `json:"message"`
}

func streamChatCompletions(
	endpoint string,
	requestBody []byte,
	authorizationHeader string,
	onDelta func(string),
) (string, int, error) {
	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return "", 0, err
	}

	req.Header.Add("Content-Type", "application/json")
	if strings.TrimSpace(authorizationHeader) != "" {
		req.Header.Add("Authorization", authorizationHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", 0, err
	}
	if resp == nil {
		return "", 0, fmt.Errorf("received nil response while streaming chat completions")
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", statusCode, fmt.Errorf("stream request failed with status code %d: %s", statusCode, strings.TrimSpace(string(responseBody)))
	}

	assistantResponse, err := parseStreamChatCompletionsResponse(resp.Body, onDelta)
	return assistantResponse, statusCode, err
}

func parseStreamChatCompletionsResponse(body io.Reader, onDelta func(string)) (string, error) {
	streamReader := bufio.NewReader(body)
	assistantResponse := strings.Builder{}
	rawResponse := strings.Builder{}
	sawDataLine := false

	for {
		line, err := streamReader.ReadString('\n')
		if line != "" {
			rawResponse.WriteString(line)
		}

		if err != nil && err != io.EOF {
			return assistantResponse.String(), err
		}

		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "data:") {
			sawDataLine = true
			payload := strings.TrimSpace(strings.TrimPrefix(trimmedLine, "data:"))
			if payload == "" {
				if err == io.EOF {
					break
				}
				continue
			}

			if payload == "[DONE]" {
				break
			}

			deltaContent, deltaErr := extractStreamDeltaContent(payload)
			if deltaErr != nil {
				return assistantResponse.String(), deltaErr
			}

			if deltaContent != "" {
				assistantResponse.WriteString(deltaContent)
				if onDelta != nil {
					onDelta(deltaContent)
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	if sawDataLine {
		return assistantResponse.String(), nil
	}

	fallbackResponse, fallbackErr := parseNonStreamFallback(rawResponse.String())
	if fallbackErr != nil {
		return assistantResponse.String(), fallbackErr
	}

	if fallbackResponse != "" && onDelta != nil {
		onDelta(fallbackResponse)
	}

	return fallbackResponse, nil
}

func extractStreamDeltaContent(payload string) (string, error) {
	chunk := streamChatCompletionsChunk{}
	if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
		return "", err
	}

	if chunk.Error != nil && strings.TrimSpace(chunk.Error.Message) != "" {
		return "", errors.New(chunk.Error.Message)
	}

	deltaBuilder := strings.Builder{}
	for _, choice := range chunk.Choices {
		if strings.TrimSpace(choice.Text) != "" {
			deltaBuilder.WriteString(choice.Text)
			continue
		}

		deltaText := extractChunkContent(choice.Delta)
		if deltaText == "" {
			deltaText = extractChunkContent(choice.Message)
		}

		if deltaText != "" {
			deltaBuilder.WriteString(deltaText)
		}
	}

	return deltaBuilder.String(), nil
}

func extractChunkContent(content any) string {
	switch typedContent := content.(type) {
	case string:
		return typedContent
	case map[string]any:
		if text, ok := typedContent["text"].(string); ok {
			return text
		}
		if nestedContent, exists := typedContent["content"]; exists {
			return extractChunkContent(nestedContent)
		}
		return ""
	case []any:
		builder := strings.Builder{}
		for _, part := range typedContent {
			builder.WriteString(extractChunkContent(part))
		}
		return builder.String()
	default:
		return ""
	}
}

func parseNonStreamFallback(rawResponse string) (string, error) {
	trimmedResponse := strings.TrimSpace(rawResponse)
	if trimmedResponse == "" {
		return "", fmt.Errorf("received empty response from stream endpoint")
	}

	chatResponse := OpenAIGPTResponse{}
	if err := json.Unmarshal([]byte(trimmedResponse), &chatResponse); err != nil {
		return "", err
	}

	if len(chatResponse.Choices) == 0 {
		return "", nil
	}

	return chatResponse.Choices[0].Message.Content, nil
}
