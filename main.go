package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"main/logging"
	"main/messages"
	"main/util"
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
	OpenAIKey       string `json:"OpenAIKey"`
	DiscordToken    string `json:"DiscordToken"`
	DevDiscordToken string `json:"DevDiscordToken"`
}

var (
	config        Config
	createdConfig bool
	devMode       bool

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},	
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
	}
)

func readConfig() {
	var bConfig []byte

	bConfig, err := os.ReadFile("config.json")
	if err != nil {
		createdConfig = true
		log.Printf("Error reading config. Creating File")
		os.WriteFile("config.json", []byte("{\"DiscordToken\":\"\",\"OpenAIKey\":\"\"}"), 0666)
		bConfig, err = os.ReadFile("config.json")
		util.HandleFatalErrors(err, "Could not read config file: config.json")
	}

	err = json.Unmarshal(bConfig, &config)
	util.HandleFatalErrors(err, "Could not parse: config_old.json")

	// Handling the case the config file has just been created
	if config.DiscordToken == "" {
		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter the Discord Token: ")
		config.DiscordToken, _ = reader.ReadString('\n')
		config.DiscordToken = strings.TrimSuffix(config.DiscordToken, "\r\n")
		log.Println("Discord Token Set to: '" + config.DiscordToken + "'")
	}

	if config.OpenAIKey == "" {
		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter the Open AI Key: ")
		config.OpenAIKey, _ = reader.ReadString('\n')
		config.OpenAIKey = strings.TrimSuffix(config.OpenAIKey, "\r\n")
		log.Println("Open AI Key Set to: '" + config.OpenAIKey + "'")
	}

}

func main() {

	createdConfig = false

	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.Parse()

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

	// Create the Discord client and add the handler to process messages
	var discord *discordgo.Session

	// TODO - Handle case where a user enters dev mode and there isnt a dev mode key
	if devMode {
		log.Println("Entering Dev Mode...")
		discord, err = discordgo.New("Bot " + config.DevDiscordToken);
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	} else {
		discord, err = discordgo.New("Bot " + config.DiscordToken)
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	}

	discord.AddHandler(messageReceive)

	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket")
		log.Fatal(err.Error())
	}

	// Used for adding slash commands
	discord.ApplicationCommandCreate(discord.State.User.ID, "", commands[0]);
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i);
		}
	})

	fmt.Println("ID: " + discord.State.User.ID);

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
	if createdConfig {
		fmt.Println("Config Updated. Saving...")
		fle, _ := json.Marshal(config)
		err := os.WriteFile("config.json", fle, 0666)
		if err != nil {
			log.Fatalf("Error Writing config_old.json")
			return
		}
	}
	logging.ShutDown()
	_ = discord.Close()
}

// dev - https://discord.com/oauth2/authorize?client_id=1009233301778743457&scope=bot&permissions=2048
// This - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2048
// https://beta.openai.com/account/usage
