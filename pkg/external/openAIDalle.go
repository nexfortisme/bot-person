package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const defaultImageOutputFormat = "png"

type openRouterImageURL struct {
	URL string `json:"url"`
}

type openRouterImageReference struct {
	Type     string             `json:"type"`
	ImageURL openRouterImageURL `json:"image_url"`
}

type openRouterImageRequest struct {
	Model           string                     `json:"model"`
	Prompt          string                     `json:"prompt"`
	N               int                        `json:"n,omitempty"`
	AspectRatio     string                     `json:"aspect_ratio,omitempty"`
	OutputFormat    string                     `json:"output_format,omitempty"`
	InputReferences []openRouterImageReference `json:"input_references,omitempty"`
}

func GetDalleResponse(prompt string) (discordgo.File, error) {
	return getImageGenerationResponse(prompt, defaultImageOutputFormat)
}

func GetDalleGifResponse(prompt string) (discordgo.File, error) {
	model := strings.ToLower(ImageGenerationModel())
	if !strings.Contains(model, "gpt-image") {
		return discordgo.File{}, errors.New("GIF generation requires a GPT image model")
	}

	return getImageGenerationResponse(prompt, "gif")
}

func getImageGenerationResponse(prompt string, outputFormat string) (discordgo.File, error) {
	openRouterResponse, err := postOpenRouterImageRequest(newImageGenerationRequest(prompt, outputFormat))
	if err != nil {
		return discordgo.File{}, err
	}

	return imageResponseToDiscordFile(openRouterResponse)
}

func GetImageEditResponse(prompt string, imageURLs []string) (discordgo.File, error) {
	urls := make([]string, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		if strings.TrimSpace(imageURL) != "" {
			urls = append(urls, imageURL)
		}
	}
	if len(urls) == 0 {
		return discordgo.File{}, errors.New("no image was provided to edit")
	}

	openRouterResponse, err := postOpenRouterImageRequest(newImageEditRequest(prompt, urls))
	if err != nil {
		return discordgo.File{}, err
	}

	return imageResponseToDiscordFile(openRouterResponse)
}

func newImageGenerationRequest(prompt string, outputFormat string) openRouterImageRequest {
	return openRouterImageRequest{
		Model:        ImageGenerationModel(),
		Prompt:       prompt,
		N:            1,
		AspectRatio:  "1:1",
		OutputFormat: outputFormat,
	}
}

func newImageEditRequest(prompt string, imageURLs []string) openRouterImageRequest {
	images := make([]openRouterImageReference, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		images = append(images, openRouterImageReference{
			Type:     "image_url",
			ImageURL: openRouterImageURL{URL: imageURL},
		})
	}

	return openRouterImageRequest{
		Model:           ImageGenerationModel(),
		Prompt:          prompt,
		N:               1,
		AspectRatio:     "1:1",
		OutputFormat:    defaultImageOutputFormat,
		InputReferences: images,
	}
}

func postOpenRouterImageRequest(requestBody openRouterImageRequest) (DalleResponse, error) {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return DalleResponse{}, errors.New("error creating request body: " + err.Error())
	}

	postRequest, err := http.NewRequest(http.MethodPost, openRouterImagesURL, bytes.NewReader(body))
	if err != nil {
		return DalleResponse{}, errors.New("POST Request Error: " + err.Error())
	}
	setOpenRouterHeaders(postRequest)

	httpResponse, err := openRouterImageClient.Do(postRequest)
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

	var openRouterResponse DalleResponse
	err = json.Unmarshal(responseBuffer, &openRouterResponse)
	if err != nil {
		return DalleResponse{}, errors.New("error parsing response: " + err.Error())
	}

	if openRouterResponse.Error != nil {
		return DalleResponse{}, errors.New("Error response from OpenRouter: " + SimplifyOpenAIError(openRouterResponse.Error.Message))
	}
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return DalleResponse{}, errors.New("API Error: " + httpResponse.Status)
	}

	return openRouterResponse, nil
}

func imageResponseToDiscordFile(openRouterResponse DalleResponse) (discordgo.File, error) {
	if len(openRouterResponse.Data) == 0 {
		return discordgo.File{}, errors.New("no image was generated")
	}

	image := openRouterResponse.Data[0]
	imageData, mediaType, err := readGeneratedImage(image)
	if err != nil {
		return discordgo.File{}, err
	}

	outputFormat := imageFormat(image, mediaType, openRouterResponse.OutputFormat)
	fileObj := &discordgo.File{
		Name:        fmt.Sprintf("%s.%s", imageFileName(openRouterResponse.Created), outputFormat),
		ContentType: imageContentType(outputFormat),
		Reader:      bytes.NewReader(imageData),
	}

	return *fileObj, nil
}

func readGeneratedImage(image DalleImages) ([]byte, string, error) {
	if strings.TrimSpace(image.B64JSON) != "" {
		imageData, err := decodeBase64Image(image.B64JSON)
		if err != nil {
			return nil, "", errors.New("Error decoding image: " + err.Error())
		}
		return imageData, image.MediaType, nil
	}

	if strings.TrimSpace(image.URL) == "" {
		return nil, "", errors.New("no image data was returned")
	}

	response, err := openRouterImageClient.Get(image.URL)
	if err != nil {
		return nil, "", errors.New("error downloading image: " + err.Error())
	}
	if response == nil {
		return nil, "", errors.New("error downloading image: empty response")
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, "", errors.New("error downloading image: " + response.Status)
	}

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, "", errors.New("error downloading image: " + err.Error())
	}
	if len(imageData) == 0 {
		return nil, "", errors.New("no image data was returned")
	}

	mediaType := image.MediaType
	if strings.TrimSpace(mediaType) == "" {
		mediaType = response.Header.Get("Content-Type")
	}
	return imageData, mediaType, nil
}

func imageFileName(created int64) string {
	if created > 0 {
		return fmt.Sprintf("bot-image-%d", created)
	}

	return fmt.Sprintf("bot-image-%s", time.Now().Format("2006-01-02-15-04-05"))
}

func imageFormat(image DalleImages, mediaType string, outputFormat string) string {
	if format := formatFromMediaType(mediaType); format != "" {
		return format
	}
	if format := formatFromMediaType(image.MediaType); format != "" {
		return format
	}
	if strings.TrimSpace(outputFormat) != "" {
		return normalizeImageFormat(outputFormat)
	}
	return defaultImageOutputFormat
}

func formatFromMediaType(mediaType string) string {
	mediaType = strings.ToLower(mediaType)
	switch {
	case strings.Contains(mediaType, "jpeg"), strings.Contains(mediaType, "jpg"):
		return "jpg"
	case strings.Contains(mediaType, "webp"):
		return "webp"
	case strings.Contains(mediaType, "gif"):
		return "gif"
	case strings.Contains(mediaType, "png"):
		return "png"
	default:
		return ""
	}
}

func normalizeImageFormat(outputFormat string) string {
	switch strings.ToLower(strings.TrimSpace(outputFormat)) {
	case "jpg", "jpeg":
		return "jpg"
	case "webp":
		return "webp"
	case "gif":
		return "gif"
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
	case "gif":
		return "image/gif"
	default:
		return "image/png"
	}
}

// ImageEditsSupported reports whether image edits can be sent to the
// configured image model. GPT Image 2.6 Sunburst accepts reference images.
func ImageEditsSupported() bool {
	return true
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
