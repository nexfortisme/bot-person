package main

import (
	"flag"
	"fmt"

	// "io"
	"log"

	"main/pkg/commands"
	"main/pkg/messages"
	"main/pkg/util"

	persistance "main/pkg/persistance"
	persistanceServices "main/pkg/persistance/services"

	loggingTypes "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// config  util.ConfigStruct
	devMode bool

	fiveMinuteTicker = time.NewTicker(5 * time.Minute)

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

func main() {
	// Step One: Read in flag variables
	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.BoolVar(&removeCommands, "removeCommands", false, "Flag for removing registered commands on shutdown")
	flag.Parse()

	// Step 2: Read in environment variables
	util.ReadEnv()

	// Step 3: Connect to the database
	databseConnection := persistance.GetDB()

	_, insertError := logging.LogEvent(loggingTypes.BOT_START, "Bot Person is starting up.", "System", "System")
	if insertError != nil {
		log.Fatalf("Error logging event: %v", insertError)
	}

	// Step 4: Declare and create the Discord Session
	var discordSession *discordgo.Session
	var err error

	if devMode {
		log.Println("Entering Dev Mode...")
		discordSession, err = discordgo.New("Bot " + util.GetDevDiscordKey())
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	} else {
		discordSession, err = discordgo.New("Bot " + util.GetDiscordKey())
		if err != nil {
			log.Fatal("Error connecting bot to server")
		}
	}

	// Step 5: Add the messageReceive handler to the discord session
	discordSession.AddHandler(messageReceive)

	// Step 6: Open the discord session
	err = discordSession.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket. " + err.Error())
	}

	// Step 7: Register the slash commands
	if removeCommands {
		removeRegisteredSlashCommands(discordSession)
	}

	if !skipCmdReg {
		registerSlashCommands(discordSession)
	}

	// Step 8: Done
	log.Println("Bot is now running")

	fmt.Println("Getting User")
	persistanceServices.GetUser("1", discordSession)

	// Pulled from the examples for discordgo, this lets the bot continue to run
	// until an interrupt is received, at which point the bot disconnects from
	// the server cleanly
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	defer databseConnection.Close()

	// This is a simple 5 minute loop originally used to save the bot statistics
	for {
		select {
		case <-fiveMinuteTicker.C:
			// saveBotStatistics()
		case <-interrupt:
			fmt.Println("Interrupt received, stopping...")
			fiveMinuteTicker.Stop()
			shutDown(discordSession)
			logging.LogEvent(loggingTypes.BOT_STOP, "Bot Person is shutting down.", "System", "System")
			return
		}
	}

}

func messageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	messages.ParseMessage(s, m)
}

func registerNewCommands(s *discordgo.Session) {
	log.Println("Registering Commands...")

	registeredCommands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}
	// Loop through new commands and check if they are already registered
	for _, newCmd := range slashCommands {
		var found bool
		for _, existingCmd := range registeredCommands {
			if newCmd.Name == existingCmd.Name {
				found = true
				break
			}
		}

		// If the command is not found among existing ones, register it
		if !found {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, "", newCmd)
			if err != nil {
				fmt.Printf("Could not create command %s: %v\n", newCmd.Name, err)
			} else {
				fmt.Printf("Successfully registered new command: %s\n", newCmd.Name)
			}
		}
	}
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

// I guess this is redundant since we no longer have to worry about
// persisting bot data to a file
func shutDown(discord *discordgo.Session) {
	log.Println("Shutting Down...")
	_ = discord.Close()
}
