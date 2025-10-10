package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func GetDalleResponse(prompt string) (discordgo.File, error) {
	httpClient := &http.Client{}

	requestDataTemplate := `{
		"model": "%s",
		"input": [
			{
				"role": "developer",
				"content": "You help with the users requests, but if there is an error, do not ask for follow up. Do not give them clarification on what they could change. Just give the reason for the error and nothing else."
			},
			{
				"role": "user",
				"content": "%s"
			}
		],
		"tools": [{"type": "image_generation"}]
	  }`
	requestData := fmt.Sprintf(requestDataTemplate, util.GetImageGenerationModel(), prompt)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", strings.NewReader(requestData))
	if err != nil {
		return discordgo.File{}, errors.New("POST Request Error: " + err.Error())
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpResponse, _ := httpClient.Do(postRequest)
	if httpResponse == nil {
		return discordgo.File{}, errors.New("API Error: " + httpResponse.Status)
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)

	var openAIResponse Response
	var imageResult *Output

	err = json.Unmarshal([]byte(string(responseBuffer)), &openAIResponse)
	if err != nil {
		return discordgo.File{}, errors.New("error Parsing Response")
	}

	// -------------------- DEBUGGING --------------------
	currentTime := time.Now().Format("2006-01-02-15-04-05")

	data, _ := json.MarshalIndent(openAIResponse, "", "  ")
	util.SaveResponseToFile(data, fmt.Sprintf("dalle-response-%s.txt", currentTime))
	// -------------------- DEBUGGING --------------------

	imageResult, err = checkResponseForErrors(openAIResponse)
	if err != nil {
		return discordgo.File{}, errors.New("Error Response from OpenAI: " + err.Error())
	}

	err = saveImageResponseToFile(imageResult.Result, openAIResponse.ID)
	if err != nil {
		return discordgo.File{}, errors.New("Error saving image response to file: " + err.Error())
	}

	reader, err := os.Open(filepath.Join("img", fmt.Sprintf("%s.jpg", openAIResponse.ID)))
	if err != nil {
		return discordgo.File{}, errors.New("error opening file: " + err.Error())
	}

	fileInfo, err := reader.Stat()
	if err != nil {
		return discordgo.File{}, errors.New("error getting file info: " + err.Error())
	}

	fileObj := &discordgo.File{
		Name:        fileInfo.Name(),
		ContentType: "image/jpg",
		Reader:      reader,
	}

	return *fileObj, nil
}

func GetDalleFollowupResponse(prompt string, previous_response_id string) (discordgo.File, error) {
	httpClient := &http.Client{}

	followUpRequestDataTemplate := `{
		"model": "%s",
		"input": [
			{
	 			"role": "developer",
	 			"content": "If the message from the user apperas to not reference the previous response, do not ask for follow up and return an error."
	 		},
	 		{
	 			"role": "user",
	 			"content": "%s"
	 		}
		],
		"previous_response_id": "%s",
		"tools": [{"type": "image_generation"}]
	}`

	requestData := fmt.Sprintf(followUpRequestDataTemplate, util.GetImageGenerationModel(), prompt, previous_response_id)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/responses", strings.NewReader(requestData))
	if err != nil {
		return discordgo.File{}, errors.New("POST Request Error: " + err.Error())
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpResponse, _ := httpClient.Do(postRequest)
	if httpResponse == nil {
		return discordgo.File{}, errors.New("API Error: " + httpResponse.Status)
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)

	var openAIResponse2 Response
	var imageResult2 *Output

	err = json.Unmarshal([]byte(string(responseBuffer)), &openAIResponse2)
	if err != nil {
		return discordgo.File{}, errors.New("error parsing response: " + err.Error())
	}

	// -------------------- DEBUGGING --------------------
	currentTime := time.Now().Format("2006-01-02-15-04-05")

	data, _ := json.MarshalIndent(openAIResponse2, "", "  ")
	util.SaveResponseToFile(data, fmt.Sprintf("dalle-followup-response-%s.txt", currentTime))
	// -------------------- DEBUGGING --------------------

	imageResult2, err = checkResponseForErrors(openAIResponse2)
	if err != nil {
		return discordgo.File{}, errors.New("Error Response from OpenAI: " + err.Error())
	}

	err = saveImageResponseToFile(imageResult2.Result, openAIResponse2.ID)
	if err != nil {
		return discordgo.File{}, errors.New("Error saving image response to file: " + err.Error())
	}

	reader, err := os.Open(filepath.Join("img", fmt.Sprintf("%s.jpg", openAIResponse2.ID)))
	if err != nil {
		return discordgo.File{}, errors.New("error opening file: " + err.Error())
	}

	fileInfo, err := reader.Stat()
	if err != nil {
		return discordgo.File{}, errors.New("error getting file info: " + err.Error())
	}

	fileObj := &discordgo.File{
		Name:        fileInfo.Name(),
		ContentType: "image/jpg",
		Reader:      reader,
	}

	return *fileObj, nil
}

func saveImageResponseToFile(responseB64String string, fileName string) error {

	path := filepath.Join("img", fmt.Sprintf("%s.jpg", removePunctuation(fileName)))

	b64 := responseB64String
	if i := strings.Index(b64, ","); i != -1 && strings.Contains(b64[:i], ";base64") {
		b64 = b64[i+1:]
	}

	// Try standard base64 first
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		// Fallback for URL-safe / unpadded strings
		data, err = base64.RawStdEncoding.DecodeString(b64)
		if err != nil {
			// Last try: URL encoding
			data, err = base64.RawURLEncoding.DecodeString(b64)
		}
	}
	if err != nil {
		fmt.Println("invalid base64: " + err.Error())
		return err
		// return discordgo.File{}, errors.New("invalid base64: " + err.Error())
	}

	// Create a new file to save the image to
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the image data to the file
	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return err
	}

	fmt.Println("Image saved to: " + path)

	return nil
}

func checkResponseForErrors(openAIResponse Response) (*Output, error) {

	var imageResult *Output

	// Types of Errors:
	// 1. Safety Error
	// 2. API Error
	// 3. No Image

	if openAIResponse.Error != nil {
		return nil, errors.New(SimplifyOpenAIError(openAIResponse.Error.Message))
	}

	if len(openAIResponse.Output) == 0 {
		return nil, errors.New(`no image was generated`)
	}

	// Checking for image generation result in output
	for _, result := range openAIResponse.Output {
		if result.Type == "image_generation_call" {
			imageResult = &result
			break
		}
	}

	// Image wasn't generated. Gave some reason why.
	if imageResult == nil {
		var completedResponse *Output

		for i := range openAIResponse.Output {
			if openAIResponse.Output[i].Status == "completed" {
				completedResponse = &openAIResponse.Output[i]
			}
		}

		if completedResponse == nil {
			return nil, errors.New(`no image was generated`)
		}

		return nil, errors.New(completedResponse.Content[0].Text)
	}

	return imageResult, nil
}

func SimplifyOpenAIError(errMsg string) string {
	msg := strings.ToLower(errMsg)

	// Safety system errors
	if strings.Contains(msg, "safety system") {
		// Extract safety violation type, e.g. safety_violations=[sexual]
		re := regexp.MustCompile(`safety_violations=\[(.*?)\]`)
		matches := re.FindStringSubmatch(msg)
		if len(matches) > 1 {
			return "Request was blocked by the safety system for " + matches[1] + " content"
		}
		return "Request was blocked by the safety system"
	}

	// Rate limiting or quota errors
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "quota") {
		return "Too many requests — please wait and try again later"
	}

	// Authentication / key issues
	if strings.Contains(msg, "invalid api key") || strings.Contains(msg, "unauthorized") {
		return "Invalid or missing API key"
	}

	// Connection or timeout
	if strings.Contains(msg, "timeout") || strings.Contains(msg, "connection refused") {
		return "Connection error — please check your network or try again later"
	}

	// Unknown / fallback
	return "An error occurred: " + errMsg
}

func removePunctuation(s string) string {
	var result strings.Builder
	for _, c := range s {
		if !strings.ContainsAny(string(c), ",.?!;:-") {
			result.WriteRune(c)
		}
	}
	return result.String()
}

func truncateString(input string) string {
	if len(input) > 50 {
		return input[:50]
	}
	return input
}
