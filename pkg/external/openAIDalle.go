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
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	openAIImageGenerationsEndpoint = "https://api.openai.com/v1/images/generations"
	openAIImageEditsEndpoint       = "https://api.openai.com/v1/images/edits"
	defaultImageSize               = "1024x1024"
	defaultImageOutputFormat       = "png"
)

type imageGenerationRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	OutputFormat   string `json:"output_format,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
}

type imageEditRequest struct {
	Model        string           `json:"model"`
	Prompt       string           `json:"prompt"`
	Images       []imageReference `json:"images"`
	N            int              `json:"n,omitempty"`
	Size         string           `json:"size,omitempty"`
	OutputFormat string           `json:"output_format,omitempty"`
}

type imageReference struct {
	ImageURL string `json:"image_url,omitempty"`
	FileID   string `json:"file_id,omitempty"`
}

func GetDalleResponse(prompt string) (discordgo.File, error) {
	openAIResponse, err := postOpenAIImageRequest(
		openAIImageGenerationsEndpoint,
		newImageGenerationRequest(prompt),
		"dalle-response",
	)
	if err != nil {
		return discordgo.File{}, err
	}

	return imageResponseToDiscordFile(openAIResponse)
}

func GetDalleFollowupResponse(prompt string, previousImageURL string) (discordgo.File, error) {
	if strings.TrimSpace(previousImageURL) == "" {
		return discordgo.File{}, errors.New("follow-up image URL is missing")
	}

	if !supportsImageEdits(util.GetImageGenerationModel()) {
		return discordgo.File{}, errors.New("follow-up image edits require a GPT image model")
	}

	openAIResponse, err := postOpenAIImageRequest(
		openAIImageEditsEndpoint,
		newImageEditRequest(prompt, previousImageURL),
		"dalle-followup-response",
	)
	if err != nil {
		return discordgo.File{}, err
	}

	return imageResponseToDiscordFile(openAIResponse)
}

func newImageGenerationRequest(prompt string) imageGenerationRequest {
	model := util.GetImageGenerationModel()
	request := imageGenerationRequest{
		Model:  model,
		Prompt: prompt,
		N:      1,
		Size:   defaultImageSize,
	}

	if usesGPTImageOutputFields(model) {
		request.OutputFormat = defaultImageOutputFormat
	} else {
		request.ResponseFormat = "b64_json"
	}

	return request
}

func newImageEditRequest(prompt string, previousImageURL string) imageEditRequest {
	return imageEditRequest{
		Model:        util.GetImageGenerationModel(),
		Prompt:       prompt,
		Images:       []imageReference{{ImageURL: previousImageURL}},
		N:            1,
		Size:         defaultImageSize,
		OutputFormat: defaultImageOutputFormat,
	}
}

func postOpenAIImageRequest(endpoint string, requestBody any, debugPrefix string) (DalleResponse, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return DalleResponse{}, errors.New("error creating request body: " + err.Error())
	}

	postRequest, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return DalleResponse{}, errors.New("POST Request Error: " + err.Error())
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpClient := &http.Client{}
	httpResponse, err := httpClient.Do(postRequest)
	if err != nil {
		return DalleResponse{}, errors.New("API Error: " + err.Error())
	}
	if httpResponse == nil {
		return DalleResponse{}, errors.New("API Error: empty response")
	}
	defer httpResponse.Body.Close()

	responseBuffer, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return DalleResponse{}, errors.New("error reading response: " + err.Error())
	}

	var openAIResponse DalleResponse
	err = json.Unmarshal(responseBuffer, &openAIResponse)
	if err != nil {
		return DalleResponse{}, errors.New("error parsing response: " + err.Error())
	}

	saveImageResponseForDebugging(openAIResponse, debugPrefix)

	if openAIResponse.Error != nil {
		return DalleResponse{}, errors.New("Error Response from OpenAI: " + SimplifyOpenAIError(openAIResponse.Error.Message))
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return DalleResponse{}, errors.New("API Error: " + httpResponse.Status)
	}

	return openAIResponse, nil
}

func imageResponseToDiscordFile(openAIResponse DalleResponse) (discordgo.File, error) {
	if len(openAIResponse.Data) == 0 {
		return discordgo.File{}, errors.New("no image was generated")
	}

	image := openAIResponse.Data[0]
	if strings.TrimSpace(image.B64JSON) == "" {
		return discordgo.File{}, errors.New("no image data was returned")
	}

	imageData, err := decodeBase64Image(image.B64JSON)
	if err != nil {
		return discordgo.File{}, errors.New("Error decoding image: " + err.Error())
	}

	outputFormat := normalizeImageFormat(openAIResponse.OutputFormat)
	fileObj := &discordgo.File{
		Name:        fmt.Sprintf("%s.%s", imageFileName(openAIResponse.Created), outputFormat),
		ContentType: imageContentType(outputFormat),
		Reader:      bytes.NewReader(imageData),
	}

	return *fileObj, nil
}

func imageFileName(created int64) string {
	if created > 0 {
		return fmt.Sprintf("openai-image-%d", created)
	}

	return fmt.Sprintf("openai-image-%s", time.Now().Format("2006-01-02-15-04-05"))
}

func normalizeImageFormat(outputFormat string) string {
	switch strings.ToLower(strings.TrimSpace(outputFormat)) {
	case "jpg", "jpeg":
		return "jpg"
	case "webp":
		return "webp"
	default:
		return "png"
	}
}

func imageContentType(outputFormat string) string {
	switch outputFormat {
	case "jpg":
		return "image/jpeg"
	case "webp":
		return "image/webp"
	default:
		return "image/png"
	}
}

func saveImageResponseForDebugging(openAIResponse DalleResponse, debugPrefix string) {
	currentTime := time.Now().Format("2006-01-02-15-04-05")
	data, _ := json.MarshalIndent(openAIResponse, "", "  ")
	util.SaveResponseToFile(data, fmt.Sprintf("%s-%s.txt", debugPrefix, currentTime))
}

func usesGPTImageOutputFields(model string) bool {
	return strings.HasPrefix(strings.ToLower(model), "gpt-image-") || strings.EqualFold(model, "chatgpt-image-latest")
}

func supportsImageEdits(model string) bool {
	return usesGPTImageOutputFields(model)
}

func decodeBase64Image(b64String string) ([]byte, error) {
	b64 := b64String
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
		return nil, errors.New("invalid base64: " + err.Error())
	}

	return data, nil
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
