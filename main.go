package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"main/logging"
	"main/messages"
	"os"
	"os/signal"
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
	// botTracking BotTracking
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

func main() {

	readConfig()
	logging.InitBotStatistics()

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

// TODO - Do this better
func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages.ParseMessage(s, m, config.OpenAIKey)
}

func shutDown(discord *discordgo.Session) {
	fmt.Println("Shutting Down")
	logging.ShutDown()
	_ = discord.Close()
}

// This - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2048
// https://beta.openai.com/account/usage
