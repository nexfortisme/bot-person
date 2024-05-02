package external

import (
	"fmt"
	"io"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"
	"main/pkg/util"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type Message struct {
    message    string
    connection *discordgo.VoiceConnection
}

var messageQueue = make(chan Message, 100) // buffer size of 100

func ElevenLabs(text string, connection *discordgo.VoiceConnection) error {
	url := "https://api.elevenlabs.io/v1/text-to-speech/onwK4e9ZLuTAKqWW03F9"

	payloadString := fmt.Sprintf("{\n  \"text\": \"%s\",\n  \"voice_settings\": {\n    \"stability\": 0.5,\n    \"similarity_boost\": 0.5\n  }\n}", text)

	payload := strings.NewReader(payloadString)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("xi-api-key", util.GetElevenLabsKey())
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	// Create a new file
	out, err := os.Create("response.mp3")
	if err != nil {
		// handle error
	}
	defer out.Close()

	// Copy the response body to the file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		// handle error
	}

	dgvoice.PlayAudioFile(connection, "response.mp3", make(chan bool))

	return nil
}

func ProcessElevenlabsMessage(message string, m *discordgo.MessageCreate ,connection *discordgo.VoiceConnection) {
	
	var messageString string
	pattern := `(?i)(?:https?:\/\/)?[\w\-\.]+\.[a-zA-Z]{2,}(\/[^\s]*)?`

	re, err := regexp.Compile(pattern)
    if err != nil {
        fmt.Println("Error compiling regex:", err)
    }

	matches := re.MatchString(message)
	re.ReplaceAllString(message, "a link")

	if(matches) {
		messageString = fmt.Sprintf("%s posted %s", m.Author.Username, message)
	} else {
		messageString = fmt.Sprintf("%s says %s", m.Author.Username, message)
	}

	logging.LogEvent(eventType.TTS_JOIN, m.Author.ID, fmt.Sprintf("Bot Said: %s", messageString), m.GuildID)

	msg := Message {
		message: messageString,
		connection: connection,
	}

	// ElevenLabs(messageString, connection)
	messageQueue <- msg
}

func ProcessQueue() {
    for {
        // Read a message from the queue
        msg := <-messageQueue

        // Process the message
        ElevenLabs(msg.message, msg.connection)
    }

}
