package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GetDalleResponse(prompt string, openAIKey string) (discordgo.File, error) {
	httpClient := &http.Client{}

	requestDataTemplate := `{
		"prompt": "%s",
		"n": 1,
		"size": "1024x1024"
	  }`
	requestData := fmt.Sprintf(requestDataTemplate, prompt)

	postRequest, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/images/generations", strings.NewReader(requestData))
	if err != nil {
		return discordgo.File{}, errors.New("POST Request Error")
	}

	postRequest.Header.Add("Content-Type", "application/json")
	postRequest.Header.Add("Authorization", "Bearer "+openAIKey)

	httpResponse, _ := httpClient.Do(postRequest)

	if httpResponse == nil {
		return discordgo.File{}, errors.New("API Error")
	}

	responseBuffer, _ := io.ReadAll(httpResponse.Body)
	var openAIResponse DalleResponse
	err = json.Unmarshal([]byte(string(responseBuffer)), &openAIResponse)
	if err != nil {
		return discordgo.File{}, errors.New("Error Parsing Response")
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(openAIResponse.Data) == 0 {
		// fmt.Println(responseBuffer)
		// fmt.Println(openAIResponse)
		// return "I'm sorry, I don't understand? (Most likely picked up by OpenAi query filter).", errors.New("API Response Error")
		return discordgo.File{}, errors.New("API Response Error. (Most Likely Picked Up By OpenAI Query Filter)")
	} else {

		err := os.MkdirAll("images", os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return discordgo.File{}, errors.New("Error creating directory")
		}

		fileName := fmt.Sprintf("images/%s.jpg", prompt)
		response, err := http.Get(openAIResponse.Data[0].URL)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()

		// Create a new file to save the image to
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// Copy the image data to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			panic(err)
		}

		reader, err := os.Open(fileName)
		fileInfo, err := reader.Stat()
		if err != nil {
			return discordgo.File{}, errors.New("Error creating file")
		}

		fileObj := &discordgo.File{
			Name:        fileInfo.Name(),
			ContentType: "image/png",
			Reader:      reader,
		}

		return *fileObj, nil
	}
}
