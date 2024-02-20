package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"

	"main/pkg/commands"
	"main/pkg/messages"
	"main/pkg/persistance"
	"main/pkg/util"

	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
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
			Description: "A command to ask the bot for a response from their infinite wisdom.",
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
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "command",
					Description: "Which command you want to get help with.",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Bot",
							Value: "bot",
						},
						{
							Name:  "Bot GPT",
							Value: "bot-gpt",
						},
						{
							Name:  "My Stats",
							Value: "my-stats",
						},
						{
							Name:  "Bot Stats",
							Value: "bot-stats",
						},
						{
							Name:  "About",
							Value: "about",
						},
						{
							Name:  "Donations",
							Value: "donations",
						},
						{
							Name:  "Images",
							Value: "images",
						},
						{
							Name:  "Balance",
							Value: "balance",
						},
						{
							Name:  "Send",
							Value: "send",
						},
						{
							Name:  "Bonus",
							Value: "bonus",
						},
						{
							Name:  "Loot Box",
							Value: "loot-box",
						},
						{
							Name:  "Broken",
							Value: "broken",
						},
						{
							Name:  "Burn",
							Value: "burn",
						},
						{
							Name:  "Stocks",
							Value: "stocks",
						},
						{
							Name:  "Portfolio",
							Value: "portfolio",
						},
						{
							Name:  "Invite",
							Value: "invite",
						},
						{
							Name:  "Save Streak",
							Value: "save-streak",
						},
						{
							Name:  "Store",
							Value: "store",
						},
					},
					Required: false,
				},
			},
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
			Name:        "loot-box",
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
		{
			Name:        "invite",
			Description: "Get an invite link to invite Bot Person to your server.",
		},
		{
			Name:        "save-streak",
			Description: "Save your streak with an save streak token or purchase one for 1/2 of your current tokens",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "action",
					Description: "Action you want to complete",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Use",
							Value: "use",
						},
						{
							Name:  "Buy",
							Value: "buy",
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "store",
			Description: "Pre-purchase your save streak tokens here. And more to come!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "item",
					Description: "The item you wish to purchase",
					Type:        discordgo.ApplicationCommandOptionString,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "Help",
							Value: "help",
						},
						{
							Name:  "Save Streak Token",
							Value: "save-streak-token",
						},
					},
					Required: true,
				},
			},
		},
		/*
			Todo:
				headsOrTails
					Bet tokens and get an RNG roll of heads or tails
				gamble
					Same as the previous gamble
				economy
					A way to see the status of the bot person economy
				leaderboard
					A way to see the top 10 users with the most tokens
				Streaks
					A way to see the top 10 users with the longest streaks
		*/
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bot":         commands.Bot,
		"bot-gpt":     commands.BotGPT,
		"my-stats":    commands.MyStats,
		"bot-stats":   commands.BotStats,
		"about":       commands.About,
		"donations":   commands.Donations,
		"help":        commands.Help,
		"image":       commands.Image,
		"balance":     commands.Balance,
		"send":        commands.Send,
		"bonus":       commands.Bonus,
		"loot-box":    commands.Lootbox,
		"broken":      commands.Broken,
		"burn":        commands.Burn,
		"stocks":      commands.Stocks,
		"portfolio":   commands.Portfolio,
		"invite":      commands.Invite,
		"save-streak": commands.SaveStreak,
		"store":       commands.Store,
	}
)

func ReadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    fmt.Printf("DB_HOST: %s\n", dbHost)
    fmt.Printf("DB_USER: %s\n", dbUser)
    fmt.Printf("DB_PASSWORD: %s\n", dbPassword)
    fmt.Printf("DB_NAME: %s\n", dbName)
}

func main() {

	// https://gobyexample.com/command-line-flags
	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.BoolVar(&removeCommands, "removeCommands", false, "Flag for removing registered commands on shutdown")
	flag.BoolVar(&disableLogging, "disableLogging", false, "Flag for disabling file logging of commands passed into bot person")
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

	if util.GetFinHubToken() == "" {
		log.Println("FinnHub Key not set, please enter your key: ")

		createdConfig = true
		reader := bufio.NewReader(os.Stdin)
		newFinnHubToken, _ := reader.ReadString('\n')
		newFinnHubToken = strings.TrimSuffix(newFinnHubToken, "\r\n")

		util.SetFinnHubToken(newFinnHubToken)

		log.Println("Finn Hub Token Set to: '" + util.GetFinHubToken() + "'")
	}

	// Adding a simple message handler
	// Mostly used for "!" commands
	discordSession.AddHandler(messageReceive)
	discordSession.AddHandler(messages.InteractionCreate)

	err = discordSession.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket. " + err.Error())
	}

	if removeCommands {
		removeRegisteredSlashCommands(discordSession)
	}

	if !skipCmdReg {
		registerSlashCommands(discordSession)
	}

	log.Println("Bot is now running")

	go listenForCommands(discordSession)

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

func listenForCommands(s *discordgo.Session) {
	fmt.Print("Enter Command:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := scanner.Text()

		// Handle the command
		if strings.HasPrefix(command, "ping") {
			fmt.Println("Pong!")
		} else if strings.HasPrefix(command, "addAdmin") {
			commandRequest := strings.Split(command, " ")

			if len(commandRequest) < 2 {
				fmt.Println("Command: addAdmin <UserId>")
				fmt.Print("Enter Command:")
				continue
			}

			createdConfig = true

			util.AddAdmin(commandRequest[1])
			fmt.Printf("Added %v to Admins\n", commandRequest[1])
		} else if strings.HasPrefix(command, "removeAdmin") {
			commandRequest := strings.Split(command, " ")

			if len(commandRequest) < 2 {
				fmt.Println("Command: removeAdmin <UserId>")
				fmt.Print("Enter Command:")
				continue
			}

			createdConfig = true
			util.RemoveAdmin(commandRequest[1])
			fmt.Printf("Removed %v from Admins\n", commandRequest[1])
		} else if strings.HasPrefix(command, "listAdmins") {
			fmt.Println("Admins: " + util.ListAdmins())
		} else if command == "quit" {
			fmt.Println("Quit Recieved, stopping...")
			shutDown(s)
			os.Exit(0)
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print("Enter Command:")
	}
}
