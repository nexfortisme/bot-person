package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"main/pkg/commands"
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
	devMode    bool
	useEnvFile bool = false

	fiveMinuteTicker = time.NewTicker(5 * time.Minute)

	removeCommands   bool
	removeOnStartup  bool
	removeOnShutdown bool

	disableLogging  bool
	disableTracking bool
	skipCmdReg      bool

	fsInterrupt bool

	createdConfig = false

	slashCommands = []*discordgo.ApplicationCommand{
		(&commands.Bot{}).ApplicationCommand(),
		(&commands.BotGPT{}).ApplicationCommand(),
		(&commands.MyStats{}).ApplicationCommand(),
		(&commands.BotStats{}).ApplicationCommand(),
		(&commands.About{}).ApplicationCommand(),
		(&commands.Donations{}).ApplicationCommand(),
		(&commands.Help{}).ApplicationCommand(),
		(&commands.Image{}).ApplicationCommand(),
		(&commands.Balance{}).ApplicationCommand(),
		(&commands.Send{}).ApplicationCommand(),
		(&commands.Bonus{}).ApplicationCommand(),
		// (&commands.Lootbox{}).ApplicationCommand(),
		(&commands.Broken{}).ApplicationCommand(),
		(&commands.Burn{}).ApplicationCommand(),
		(&commands.Invite{}).ApplicationCommand(),
		// (&commands.HSRCode{}).ApplicationCommand(),
		(&commands.Search{}).ApplicationCommand(),
		(&commands.Set{}).ApplicationCommand(),
		// (&commands.Testing{}).ApplicationCommand(),
		(&commands.Slop{}).ApplicationCommand(),
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"bot":       (&commands.Bot{}).Execute,
		"bot-gpt":   (&commands.BotGPT{}).Execute,
		"my-stats":  (&commands.MyStats{}).Execute,
		"bot-stats": (&commands.BotStats{}).Execute,
		"about":     (&commands.About{}).Execute,
		"donations": (&commands.Donations{}).Execute,
		"help":      (&commands.Help{}).Execute,
		"image":     (&commands.Image{}).Execute,
		"balance":   (&commands.Balance{}).Execute,
		"send":      (&commands.Send{}).Execute,
		"bonus":     (&commands.Bonus{}).Execute,
		// "loot-box":  (&commands.Lootbox{}).Execute,
		"broken":    (&commands.Broken{}).Execute,
		"burn":      (&commands.Burn{}).Execute,
		"invite":    (&commands.Invite{}).Execute,
		// "hsr-code":  (&commands.HSRCode{}).Execute,
		"search":    (&commands.Search{}).Execute,
		"set":       (&commands.Set{}).Execute,
		// "testing":   (&commands.Testing{}).Execute,
		"slop":      (&commands.Slop{}).Execute,
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
	flag.BoolVar(&useEnvFile, "useEnvFile", true, "Flag for using the .env file")
	flag.Parse()

	// Step 2: Read in environment variables
	util.ReadEnv(useEnvFile, devMode)

	return;

	// Step 3: Connect to the database
	databseConnection := persistance.GetDB()

	// Step 4: Declare and create the Discord Session
	var discordSession *discordgo.Session
	var err error

	var inDevMode = devMode || os.Getenv("DEV_MODE") == "true"

	// log.Printf("DEV_MODE: %v", os.Getenv("DEV_MODE"))
	// log.Printf("devMode: %v", devMode)

	// Checking to see if we are in dev mode
	if inDevMode {
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
		log.Fatal("Error opening bot websocket. Error: " + err.Error())
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
