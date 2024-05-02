package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"main/pkg/commands"
	"main/pkg/external"
	"main/pkg/handlers"
	"main/pkg/logging"
	"main/pkg/messages"
	"main/pkg/persistance"
	"main/pkg/util"

	eventType "main/pkg/logging/enums"
	state "main/pkg/state/services"

	"os"
	"os/signal"
	"strings"
	"syscall"

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
							Name:  "Invite",
							Value: "invite",
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
			Description: "Spend 5 tokens to get an RNG box",
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
			Name:        "invite",
			Description: "Get an invite link to invite Bot Person to your server.",
		},
		{
			Name:        "hsr-code",
			Description: "Get the Honkai Star Rail gift code url from a code.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "prompt",
					Description: "The code to be entered",
					Required:    true,
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
		"loot-box":  commands.Lootbox,
		"broken":    commands.Broken,
		"burn":      commands.Burn,
		"invite":    commands.Invite,
		"hsr-code":  commands.HSRCode,
	}

	applicationCommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"reset_streak_button": handlers.ResetStreakButton,
		"save_streak_button":  handlers.SaveStreakButton,
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

	logging.LogEvent(eventType.BOT_START, "SYSTEM", "Bot Person is starting up.", "SYSTEM")

	// Step 5: Add the messageReceive handler to the discord session
	discordSession.AddHandler(messageReceive)

	// Step 6: Open the discord session
	err = discordSession.Open()
	if err != nil {
		log.Fatal("Error opening bot websocket. " + err.Error())
	}

	// Step 6.5: Set the session in the state
	state.SetDiscordSession(discordSession)

	// Step 7: Register the slash commands
	if removeCommands {
		removeRegisteredSlashCommands(discordSession)
	}

	if !skipCmdReg {
		registerSlashCommands(discordSession)
	}

	// Starting queue processing for tts messages
	go external.ProcessQueue()

	// Step 8: Done
	log.Println("Bot is now running")

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
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
			break
		case discordgo.InteractionMessageComponent:
			commandNameSplit := strings.Split(i.MessageComponentData().CustomID, ":")
			if h, ok := applicationCommandHandlers[commandNameSplit[0]]; ok {
				h(s, i)
			}
			break
		default:
			log.Printf("Unknown interaction type: %v", i.Type)
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
	logging.LogEvent(eventType.BOT_STOP, "SYSTEM", "Bot Person is shutting down.", "SYSTEM")
	log.Println("Shutting Down...")
	_ = discord.Close()
}
