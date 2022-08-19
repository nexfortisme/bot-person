package messages

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/logging"
	"main/util"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	OpenAIKey    string `json:"OpenAIKey"`
	DiscordToken string `json:"DiscordToken"`
}

type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIChoice struct {
	Text   string `json:"text"`
	Index  int    `json:"index"`
	Reason string `json:"finish_reason"`
}

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, openAIKey string) {

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	var incomingMessage string

	if !strings.HasPrefix(m.Message.Content, "!") {
		incomingMessage = strings.ToLower(m.Message.Content)
	} else {
		incomingMessage = m.Message.Content
	}

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {
		logging.IncrementTracker(2, m, s)
		log.Printf("Bot Person > I'm Sorry")
		// Have this have a list of responses that it pulls from randomly
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm Sorry.")
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "!botStats") {
		logging.IncrementTracker(1, m, s)
		logging.GetBotStats(s, m)
	} else if strings.HasPrefix(incomingMessage, "!myStats") {
		logging.GetUserStats(s, m)
		logging.IncrementTracker(1, m, s)
	}

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	msg := util.ReplaceIDsWithNames(m, s)

	// logging.LogIncomingMessage(s, m, msg)

	logging.IncrementTracker(1, m, s)
	respTxt := getOpenAIResponse(msg, openAIKey)

	log.Printf("Bot Person > %s \n", respTxt)
	_, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	util.HandleErrors(err)

}

func ParseSlashCommand(s *discordgo.Session, prompt string, openAIKey string) string {
	respTxt := getOpenAIResponse(prompt, openAIKey)
	respTxt = "Request: " + prompt + " " + respTxt
	log.Printf("Bot Person > %s \n", respTxt)
	return respTxt
}

func getOpenAIResponse(prompt string, openAIKey string) string {
	client := &http.Client{}

	dataTemplate := `{
		"model": "text-davinci-002",
		"prompt": "%s",
		"temperature": 0.7,
		"max_tokens": 256,
		"top_p": 1,
		"frequency_penalty": 0,
		"presence_penalty": 0
	  }`
	data := fmt.Sprintf(dataTemplate, prompt)

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", strings.NewReader(data))
	if err != nil {
		logging.LogError("Error creating POST request")
		// log.Fatalf("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+openAIKey)

	resp, _ := client.Do(req)

	if resp == nil {
		return "Error Contacting OpenAI API. Please Try Again Later."
	}

	buf, _ := ioutil.ReadAll(resp.Body)
	var rspOAI OpenAIResponse
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	err = json.Unmarshal([]byte(string(buf)), &rspOAI)
	if err != nil {
		return ""
	}

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Text
	}
}

// Checking for the !bot keyword
// TODO - Replace with slash commands
func mentionsKeyphrase(m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, "!bot") && m.Content != "!botStats"
}

// Determine if the bot's ID is in the list of users mentioned
func mentionsBot(mentions []*discordgo.User, id string) bool {
	for _, user := range mentions {
		if user.ID == id {
			return true
		}
	}
	return false
}
