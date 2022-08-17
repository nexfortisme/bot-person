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
	config         Config
	devMode        bool
	removeCommands bool

	createdConfig         = false
	integerOptionMinValue = 1.0

	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name:        "bot",
			Description: "General bot command",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "Propmt to send to the bot",
					Required:    true,
				},
			},
		},
		// {
		// 	Name:        "my-stats",
		// 	Description: "Get yout tracking data",
		// },
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "eat shit",
				},
			})
		},
		"bot": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var msg string

			if option, ok := optionMap["prompt"]; ok {
				fmt.Println("Prompt: " + option.StringValue())

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					// Ignore type for now, they will be discussed in "responses"
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						// Flags:   uint64(discordgo.MessageFlagsEphemeral),
						Content: "Thinking...",
					},
				})

				msg = messages.ParseSlashCommand(s, option.StringValue(), config.OpenAIKey)
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: msg,
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

		},
		"my-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.

			msg := logging.SlashGetUserStats(s, i)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
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

	// https://gobyexample.com/command-line-flags
	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.BoolVar(&removeCommands, "removeCommands", false, "Flag for removing registered commands on shutdown")
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

		if config.DevDiscordToken == "" {
			createdConfig = true
			reader := bufio.NewReader(os.Stdin)
			log.Print("Please Enter the Dev Discord Token: ")
			config.DevDiscordToken, _ = reader.ReadString('\n')
			config.DevDiscordToken = strings.TrimSuffix(config.DevDiscordToken, "\r\n")
			log.Println("Dev Discord Token Set to: '" + config.DevDiscordToken + "'")
		}

		discord, err = discordgo.New("Bot " + config.DevDiscordToken)
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

	registerCommands(discord)
	log.Println("Bot is now running")

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

func registerCommands(s *discordgo.Session) {
	// Used for adding slash commands
	// Add the command and then add the handler for that command
	// https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func removeRegisteredCommands(s *discordgo.Session) {
	log.Println("Removing Commands...")

	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func writeConfig() {
	log.Println("Config Updated. Writing...")

	fle, _ := json.Marshal(config)
	err := os.WriteFile("config.json", fle, 0666)
	if err != nil {
		log.Fatalf("Error Writing config.json")
		return
	}
}

func shutDown(discord *discordgo.Session) {
	log.Println("Shutting Down...")

	if createdConfig {
		writeConfig()
	}

	if removeCommands {
		removeRegisteredCommands(discord)
	}

	logging.ShutDown()
	_ = discord.Close()
}

// Dev - https://discord.com/oauth2/authorize?client_id=1009233301778743457&scope=bot&permissions=2048
// Prod - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2147485696
// https://beta.openai.com/account/usage
