package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	external "main/pkg/external/models"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
)

func GetDalleResponse(prompt string) (discordgo.File, error) {
	httpClient := &http.Client{}

	requestDataTemplate := `{
		"model": "dall-e-3",
		"prompt": "%s",
		"n": 1,
		"size": "1024x1024",
		"quality": "hd"
	  }`
	requestData := fmt.Sprintf(requestDataTemplate, prompt)

	logging.LogEvent(loggingType.EXTERNAL_DALLE_REQUEST, requestData, "System", "System", nil)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/images/generations", strings.NewReader(requestData))
	if err != nil {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "Error Creating OpenAI Request", "System", "System", nil)
		return discordgo.File{}, errors.New("POST Request Error")
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+util.GetOpenAIKey())

	httpResponse, _ := httpClient.Do(postRequest)
	if httpResponse == nil {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "Error Contacting OpenAI API", "System", "System", nil)
		return discordgo.File{}, errors.New("API Error")
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)
	var openAIResponse external.DalleResponse
	err = json.Unmarshal([]byte(string(responseBuffer)), &openAIResponse)
	if err != nil {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "Error Parsing OpenAI DALLE Response", "System", "System", nil)
		return discordgo.File{}, errors.New("error Parsing Response")
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(openAIResponse.Data) == 0 {
		logging.LogEvent(loggingType.EXTERNAL_API_ERROR, "OpenAI Response Error", "System", "System", nil)
		return discordgo.File{}, errors.New("API Response Error. (Most Likely Picked Up By OpenAI Query Filter)")
	} else {

		err = createDirectoryIfNotExists("img")
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return discordgo.File{}, errors.New("error creating directory")
		}

		path := filepath.Join("img", fmt.Sprintf("%s.jpg", TruncateString(removePunctuation(prompt))))

		response, err := http.Get(openAIResponse.Data[0].URL)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		// Create a new file to save the image to
		file, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Copy the image data to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			panic(err)
		}

		// Open the file to send it to Discord
		reader, err := os.Open(path)
		if err != nil {
			return discordgo.File{}, errors.New("error opening file")
		}

		// Get the file info
		fileInfo, err := reader.Stat()
		if err != nil {
			return discordgo.File{}, errors.New("error creating file")
		}

		// Create a new Discord file object to send
		fileObj := &discordgo.File{
			Name:        fileInfo.Name(),
			ContentType: "image/jpg",
			Reader:      reader,
		}

		responseJSON, _ := json.Marshal(openAIResponse.Data)
		logging.LogEvent(loggingType.EXTERNAL_DALLE_RESPONSE, string(responseJSON), "System", "System", nil)

		return *fileObj, nil
	}
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

func createDirectoryIfNotExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func TruncateString(input string) string {
	if len(input) > 50 {
		return input[:50]
	}
	return input
}
