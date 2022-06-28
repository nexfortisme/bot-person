package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// The outer structure of the response from OpenAI
type OpenAIResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
}

// The inter structure of the response from OpenAI, this
// contains zero or more completions based on the provided
// prompt
type OpenAIChoice struct {
	Text   string `json:"text"`
	Index  int    `json:"index"`
	Reason string `json:"finish_reason"`
}

type Config struct {
	OpenAIKey    string `json:"OpenAIKey"`
	DiscordToken string `json:"DiscordToken"`
}

var (
	config Config
)

func main() {

	// Parse the config
	bConfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Could not read config file: config.json")
	}

	err = json.Unmarshal(bConfig, &config)
	if err != nil {
		log.Fatalf("Could not parse config file")
	}

	// Create the Discord client and add the handler
	// to process messages
	discord, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		log.Fatal("Error connecting bot to server")
	}

	discord.AddHandler(messageReceive)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket")
		log.Fatal(err.Error())
	}

	fmt.Println("Bot is now running")

	// Pulled from the examples for discordgo, this lets the bot continue to run
	// until an interrupt is received, at which point the bot disconnects from
	// the server cleanly
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {

	// The bot should ignore messages from itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) {
		return
	}

	// Remove the initial mention of the bot
	toReplace := fmt.Sprintf("<@%s> ", id)
	requestUser := m.Author.Username
	msg := strings.Replace(m.Message.Content, toReplace, "", 1)
	msg = replaceMentionsWithNames(m.Mentions, msg)

	log.Printf(" %s (%s) < %s\n", requestUser, msg)

	respTxt := formulateResponse(msg)

	log.Printf("> %s\n", respTxt)
	s.ChannelMessageSend(m.ChannelID, respTxt)
}

func formulateResponse(prompt string) string {
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

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", strings.NewReader((data)))
	if err != nil {
		log.Fatalf("Error creating POST request")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.OpenAIKey)

	resp, _ := client.Do(req)

	buf, _ := ioutil.ReadAll(resp.Body)
	var rspOAI OpenAIResponse
	// TODO: This could contain an error from OpenAI (rate limit, server issue, etc)
	// need to add proper error handling
	json.Unmarshal([]byte(string(buf)), &rspOAI)

	// It's possible that OpenAI returns no response, so
	// fallback to a default one
	if len(rspOAI.Choices) == 0 {
		return "I'm sorry, I don't understand?"
	} else {
		return rspOAI.Choices[0].Text
	}
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

// The message string that the bot receives reads mentions of other users as
// an ID in the form of "<@000000000000>", instead iterate over each mention and
// replace the ID with the user's username
func replaceMentionsWithNames(mentions []*discordgo.User, message string) string {
	retStr := strings.Clone(message)
	for _, mention := range mentions {
		idStr := fmt.Sprintf("<@%s>", mention.ID)
		retStr = strings.ReplaceAll(retStr, idStr, mention.Username)
	}
	return retStr
}

// This - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2048
// https://beta.openai.com/account/usage
