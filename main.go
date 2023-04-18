package main

import (
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
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/r3labs/diff/v3"
)

var (
	// config  util.ConfigStruct
	devMode bool

	removeCommands bool

	skipCmdReg bool

	createdConfig         = false
	integerOptionMinValue = 0.1

	slashCommands = []*discordgo.ApplicationCommand{
		{
			Name:        "bot",
			Description: "Interact with OpenAI's GPT Chat Models and see what out future AI overlords have to say.",
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
			Name:        "about",
			Description: "Get information about Bot Person.",
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
							Name:  "My Stats",
							Value: "my-stats",
						},
						{
							Name:  "About",
							Value: "about",
						},
						{
							Name:  "Image",
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
							Name:  "Burn",
							Value: "burn",
						},
						//{
						//	Name:  "Invite",
						//	Value: "invite",
						//},
						{
							Name:  "Save Streak",
							Value: "save-streak",
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
					Name:        "recipient",
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
		//{
		//	Name:        "invite",
		//	Description: "Get an invite link to invite Bot Person to your server.",
		//},
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
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bot":      commands.BotGPT,
		"my-stats": commands.MyStats,
		"about":    commands.About,
		"help":     commands.Help,
		"image":    commands.Image,
		"balance":  commands.Balance,
		"send":     commands.Send,
		"bonus":    commands.Bonus,
		"loot-box": commands.Lootbox,
		"burn":     commands.Burn,
		//"invite":      commands.Invite,
		"save-streak": commands.SaveStreak,
	}
)

func main() {

	// https://gobyexample.com/command-line-flags
	flag.BoolVar(&devMode, "dev", false, "Flag for starting the bot in dev mode")
	flag.BoolVar(&removeCommands, "removeCommands", false, "Flag for removing registered commands on shutdown")
	flag.BoolVar(&skipCmdReg, "skipCmdReg", false, "Flag for disabling registering of commands on startup")
	flag.Parse()

	util.ReadConfig()
	persistance.ReadBotStatistics()

	fiveMinuteTicker := time.NewTicker(5 * time.Minute)

	logFile, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	defer logFile.Close()

	// This makes it print to both the console and to a file
	log.SetOutput(multiWriter)

	// Create the Discord client and add the handler to process messages
	var discordSession *discordgo.Session

	if devMode {
		log.Println("Entering Dev Mode...")

		if util.GetDevDiscordToken() == "" {
			createdConfig = true
			var devDiscordToken string

			util.ReadAPIKey(&devDiscordToken, "Dev Discord Token")
			util.SetDevDiscordToken(devDiscordToken)
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

	// Adding a simple message handler
	// Mostly used for "!" commands
	discordSession.AddHandler(messageReceive)

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
		persistance.SaveBotStatistics()
		persistance.ReadBotStatistics()
	} else {
		log.Println("No Changes to Bot Statistics. Skipping Save...")
	}

}
