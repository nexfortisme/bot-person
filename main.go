package main

import (
	"bufio"
	"encoding/json"
	"flag"
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
	config          Config
	devMode         bool
	removeCommands  bool
	disableLogging  bool
	disableTracking bool

	createdConfig         = false
	integerOptionMinValue = 1.0

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "test",
			Description: "A simple test commnand for the bot.",
		},
		{
			Name:        "bot",
			Description: "A command to ask the bot for a reposne from their infinite wisdom.",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "The actual prompt that the bot will ponder on.",
					Required:    true,
				},
			},
		},
		{
			Name:        "my-stats",
			Description: "Get yout tracking data.",
		},
		{
			Name:        "bot-stats",
			Description: "Get global stats for the bot.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			logging.IncrementTracker(0, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "eat shit",
				},
			})
		},
		// TODO - Handle logging of the incoming request by the user
		"bot": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var msg string

			// Pulling the propt out of the optionsMap
			if option, ok := optionMap["prompt"]; ok {

				// Generating the response
				placeholder := "Thinking about: " + option.StringValue()

				// Immediately responding in the 3 second window before the interaciton times out
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: placeholder,
					},
				})

				// Going out to make the OpenAI call to get the proper response
				msg = messages.ParseSlashCommand(s, option.StringValue(), config.OpenAIKey)

				// Incrementint interaciton counter
				logging.IncrementTracker(1, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: msg,
				})
				if err != nil {
					// Not 100% sure this is the approach I want to take with handling errors from the API
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

		},
		"my-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := logging.SlashGetUserStats(s, i)
			logging.IncrementTracker(0, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"bot-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			msg := logging.SlashGetBotStats(s)
			logging.IncrementTracker(0, i.Interaction.Member.User.ID, i.Interaction.Member.User.Username)
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
	util.HandleFatalErrors(err, "Could not parse: config.json")

	// Handling the case the config file has just been created
	if config.DiscordToken == "" {
		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		log.Print("Please Enter the Discord Token: ")
		config.DiscordToken, _ = reader.ReadString('\n')
		config.DiscordToken = strings.TrimSuffix(config.DiscordToken, "\r\n")
		log.Println("Discord Token Set to: '" + config.DiscordToken + "'")
	}

	// TODO - Check to see if the user doesn't type in a command
	// If they don't, ask them if they wish to continue without OpenAI responses
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
	flag.BoolVar(&disableLogging, "diableLogging", false, "Flag for disabling file logging of commands passed into bot person")
	flag.BoolVar(&disableTracking, "disableTracking", false, "Flag for disabling tracking of user interactions and bad bot messages")
	flag.Parse()

	readConfig()
	logging.InitBotStatistics()

	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	if !disableLogging {
		mw := io.MultiWriter(os.Stdout, f)
		defer f.Close()

		// This makes it print to both the console and to a file
		log.SetOutput(mw)
	}

	// Create the Discord client and add the handler to process messages
	var discordSession *discordgo.Session

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

		discordSession, err = discordgo.New("Bot " + config.DevDiscordToken)
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	} else {
		discordSession, err = discordgo.New("Bot " + config.DiscordToken)
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	}

	discordSession.AddHandler(messageReceive)

	err = discordSession.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket")
		log.Fatal(err.Error())
	}

	registerSlashCommands(discordSession)
	log.Println("Bot is now running")

	// Pulled from the examples for discordgo, this lets the bot continue to run
	// until an interrupt is received, at which point the bot disconnects from
	// the server cleanly
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	shutDown(discordSession)
}

// TODO - Do this better
func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages.ParseMessage(s, m, config.OpenAIKey)
}

func registerSlashCommands(s *discordgo.Session) {
	log.Println("Registering Commands...")
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

func removeRegisteredSlashCommands(s *discordgo.Session) {
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
		removeRegisteredSlashCommands(discord)
	}

	logging.ShutDown()
	_ = discord.Close()
}

// Dev - https://discord.com/oauth2/authorize?client_id=1009233301778743457&scope=bot&permissions=2048
// Prod - https://discord.com/oauth2/authorize?client_id=225979639657398272&scope=bot&permissions=2147485696
// https://beta.openai.com/account/usage
