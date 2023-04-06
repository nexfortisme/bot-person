package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"main/commands"
	"main/messages"
	"main/persistance"
	"main/util"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/r3labs/diff/v3"
)

var (
	// config  util.ConfigStruct
	devMode bool

	removeCommands   bool
	removeOnStartup  bool
	removeOnShutdown bool

	disableLogging  bool
	disableTracking bool
	skipCmdReg      bool

	fsInterrupt bool

	createdConfig         = false
	integerOptionMinValue = 0.1

	slashCommands = []*discordgo.ApplicationCommand{
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
			Name:        "bot-gpt",
			Description: "Interact with OpenAI's GPT-4 API and see what out future AI overlords have to say.",
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
			Description: "Get usage stats.",
		},
		{
			Name:        "bot-stats",
			Description: "Get global usage stats.",
		},
		{
			Name:        "about",
			Description: "Get information about Bot Person.",
		},
		{
			Name:        "donations",
			Description: "List of the people who contributed to Bot Person's on-going service.",
		},
		{
			Name:        "help",
			Description: "List of commands to use with Bot Person.",
		},
		{
			Name:        "image",
			Description: "Ask Bot Person to generate an image for you. Costs 1 Token per image",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "The actual prompt that Bot Person will generate an image from.",
					Required:    true,
				},
				// TODO - Implement requesting multiple images at once
				// {
				// 	Type:        discordgo.ApplicationCommandOptionInteger,
				// 	Name:        "number",
				// 	Description: "The number of image you want Bot Person to generate. Cost = # of images generated",
				// 	MinValue:    &integerOptionMinValue,
				// 	MaxValue:    10,
				// 	Required:    false,
				// },
			},
		},
		{
			Name:        "balance",
			Description: "Check your balance or the balance of another user.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "The person you want to check the balance of.",
					Required:    false,
				},
			},
		},
		{
			Name:        "send",
			Description: "Send tokens to another user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "recepient",
					Description: "The person you want to send tokens to.",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of tokens you want to send.",
					MinValue:    &integerOptionMinValue,
					Required:    true,
				},
			},
		},
		{
			Name:        "bonus",
			Description: "Use this command every 24 hours for a small bundle of tokens",
		},
		{
			Name:        "lootbox",
			Description: "Spend 2.5 tokens to get an RNG box",
		},
		{
			Name:        "broken",
			Description: "Get more information if something about bot person is broken",
		},
		{
			Name:        "burn",
			Description: "A way, for whatever reason, you can burn tokens.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionNumber,
					Name:        "amount",
					Description: "The amount of tokens you want to send.",
					MinValue:    &integerOptionMinValue,
					Required:    true,
				},
			},
		},
		{
			Name:        "stocks",
			Description: "Buy and sell fake stocks with bot person tokens.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "action",
					Description: "Action you want to complete",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Buy",
							Value: "buy",
						},
						{
							Name:  "Sell",
							Value: "sell",
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "stock",
					Description: "The symbol of the stock you want to buy or sell",
					Required:    true,
				},
				{
					Name:        "quantity",
					Description: "Number of stocks you want to buy or sell",
					Type:        discordgo.ApplicationCommandOptionNumber,
					MinValue:    &integerOptionMinValue,
					Required:    true,
				},
			},
		},
		{
			Name:        "portfolio",
			Description: "View your portfolio of stocks.",
		},
		/*
			Todo:
				headsOrTails
					Bet tokens and get an RNG roll of heads or tails
				gamble
					Same as the previous gamble
				economy
					A way to see the status of the bot person economy
				invite
					Generate an invite link for the bot that is specific to whatever token is being used for the bot
		*/
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bot":       commands.Bot,
		"bot-gpt":   commands.BotGPT,
		"my-stats":  commands.MyStats,
		"bot-stats": commands.BotStats,
		"about":     commands.About,
		"donations": commands.Donations,
		"help":      commands.Help,
		"image":     commands.Image,
		"balance":   commands.Balance,
		"send":      commands.Send,
		"bonus":     commands.Bonus,
		"lootbox":   commands.Lootbox,
		"broken":    commands.Broken,
		"burn":      commands.Burn,
		"stocks":    commands.Stocks,
		"portfolio": commands.Portfolio,
	}
)

func main() {

	// https://gobyexample.com/command-line-flags
	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.BoolVar(&removeCommands, "removeCommands", false, "Flag for removing registered commands on shutdown")
	flag.BoolVar(&disableLogging, "diableLogging", false, "Flag for disabling file logging of commands passed into bot person")
	flag.BoolVar(&disableTracking, "disableTracking", false, "Flag for disabling tracking of user interactions and bad bot messages")
	flag.BoolVar(&skipCmdReg, "skipCmdReg", false, "Flag for disabling registering of commands on startup")
	flag.Parse()

	util.ReadConfig()
	persistance.ReadBotStatistics()

	fiveMinuteTicker := time.NewTicker(5 * time.Minute)

	logFile, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	if !disableLogging {
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		defer logFile.Close()

		// This makes it print to both the console and to a file
		log.SetOutput(multiWriter)
	}

	// Create the Discord client and add the handler to process messages
	var discordSession *discordgo.Session

	if devMode {
		log.Println("Entering Dev Mode...")

		if util.GetDevDiscordToken() == "" {
			createdConfig = true
			reader := bufio.NewReader(os.Stdin)
			log.Print("Please Enter the Dev Discord Token: ")

			newDevToken, _ := reader.ReadString('\n')
			newDevToken = strings.TrimSuffix(newDevToken, "\r\n")

			util.SetDevDiscordToken(newDevToken)

			// config.DevDiscordToken, _ = reader.ReadString('\n')
			// config.DevDiscordToken = strings.TrimSuffix(config.DevDiscordToken, "\r\n")
			log.Println("Dev Discord Token Set to: '" + util.GetDevDiscordToken() + "'")
		}

		discordSession, err = discordgo.New("Bot " + util.GetDevDiscordToken())
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	} else {
		discordSession, err = discordgo.New("Bot " + util.GetDiscordToken())
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

	if removeCommands {
		removeRegisteredSlashCommands(discordSession)
	}

	if !skipCmdReg {
		registerSlashCommands(discordSession)
	}

	log.Println("Bot is now running")

	// Pulled from the examples for discordgo, this lets the bot continue to run
	// until an interrupt is received, at which point the bot disconnects from
	// the server cleanly
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case <-fiveMinuteTicker.C:
			saveBotStatistics()
		case <-interrupt:
			fmt.Println("Interrupt received, stopping...")
			fiveMinuteTicker.Stop()
			shutDown(discordSession)
			return
		}
	}

}

func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages.ParseMessage(s, m)
}

func registerSlashCommands(s *discordgo.Session) {
	log.Println("Registering Commands...")
	// Used for adding slash commands
	// Add the command and then add the handler for that command
	// https://github.com/bwmarrin/discordgo/blob/master/examples/slash_commands/main.go
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slashCommands))
	for i, v := range slashCommands {
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

func shutDown(discord *discordgo.Session) {
	log.Println("Shutting Down...")

	if createdConfig {
		util.WriteConfig()
	}

	persistance.SaveBotStatistics()
	_ = discord.Close()
}

// This is a function that will save the current contents of the bot statistics
func saveBotStatistics() {

	changeLog, _ := diff.Diff(persistance.GetTempTracking(), persistance.GetBotTracking())

	if len(changeLog) > 0 {
		log.Println("Saving Bot Statistics...")
		fsInterrupt = true // TODO - Implement interrupt checking for when a user may be doing something while the bot is saving
		persistance.SaveBotStatistics()
		persistance.ReadBotStatistics()
		fsInterrupt = false
	} else {
		log.Println("No Changes to Bot Statistics. Skipping Save...")
	}

}
