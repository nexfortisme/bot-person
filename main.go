package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

type BotTracking struct {
	BadBotCount  int `json:"BadBotCount"`
	MessageCount int `json:"MessageCount"`
}

var (
	config Config
)

var (
	botTracking BotTracking
)

func readConfig() {
	// Parse the config
	bConfig, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Could not read config file: config.json")
	}

	err = json.Unmarshal(bConfig, &config)
	if err != nil {
		log.Fatalf("Could not parse: config.json")
	}
}

func initBotStatistics() {
	var trackingFile []byte

	trackingFile, err := ioutil.ReadFile("botTracking.json")
	if err != nil {

		log.Printf("Error Reading botTracking. Creating File")
		ioutil.WriteFile("botTracking.json", []byte("{\"BadBotCount\":0,\"MessageCount\":0}"), 0666)

		trackingFile, err = ioutil.ReadFile("botTracking.json")
		if err != nil {
			log.Fatalf("Could not read config file: botTracking.json")
		}

	}

	err = json.Unmarshal(trackingFile, &botTracking)
	if err != nil {
		log.Fatalf("Could not parse: botTracking.json")
	}
}

func main() {

	readConfig()
	initBotStatistics()

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, f)
	defer f.Close()

	// This makes it print to both the console and to a file
	log.SetOutput(mw)

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
	shutDown(discord)

}

func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {

	var mMessage = ""
	if strings.HasPrefix(m.Message.Content, "!") {
		mMessage = m.Message.Content
	} else {
		mMessage = strings.ToLower(m.Message.Content)
	}

	// The bot should ignore messages from itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(mMessage, "bad bot") {
		logIncomingMessage(s, m, mMessage)
		incrementTracker(2)
		log.Printf("Bot Person > I'm Sorry")
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm Sorry.")
		if err != nil {
			return
		}
	} else if strings.HasPrefix(mMessage, "!badCount") {
		logIncomingMessage(s, m, mMessage)
		incrementTracker(1)
		ret := "Bad Bot Count: " + strconv.Itoa(botTracking.BadBotCount)
		_, err := s.ChannelMessageSend(m.ChannelID, ret)
		if err != nil {
			return
		}
		ret = "Bot Person > " + ret;
		log.Printf(ret)
	}

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) {
		return
	}

	// Remove the initial mention of the bot
	toReplace := fmt.Sprintf("<@%s> ", id)
	msg := strings.Replace(m.Message.Content, toReplace, "", 1)
	msg = replaceMentionsWithNames(m.Mentions, msg)

	logIncomingMessage(s, m, msg);

	respTxt := formulateResponse(msg)

	incrementTracker(1)

	log.Printf("Bot Person > %s \n", respTxt)
	_, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	if err != nil {
		return
	}
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

	req, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/completions", strings.NewReader(data))
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

// Determine if the bot's ID is in the list of users mentioned
func mentionsBot(mentions []*discordgo.User, id string) bool {
	for _, user := range mentions {
		if user.ID == id {
			return true
		}
	}
	return false
}

func mentionsKeyphrase(m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, "!bot")
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

func incrementTracker(flag int) {
	if flag == 1 {
		botTracking.MessageCount++
	} else {
		botTracking.MessageCount++
		botTracking.BadBotCount++
	}
}

func shutDown(discord *discordgo.Session) {
	fmt.Println("Shutting Down")
	fle, _ := json.Marshal(botTracking)
	ioutil.WriteFile("botTracking.json", fle, 0666)
	_ = discord.Close()
}

func logIncomingMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	requestUser := m.Author.Username
	rGuild, _ := s.State.Guild(m.GuildID)
	rGuildName := rGuild.Name

	log.Printf(" %s (%s) < %s\n", requestUser, rGuildName, message)
}

// This - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2048
// https://beta.openai.com/account/usage
