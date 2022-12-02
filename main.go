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
	"main/persistance"
	"main/util"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
	disableCmdReg   bool

	createdConfig         = false
	integerOptionMinValue = 0.1

	commands = []*discordgo.ApplicationCommand{
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
			Name: "balance",
			// TODO - Add flag for users to opt out of others being able to check their balance
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
			Description: "A command to ask the bot for a reposne from their infinite wisdom.",
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
		// {
		// 	Name:        "gamba",
		// 	Description: "Try your luck and see if you can win some extra Image Tokens.",
		// 	Options: []*discordgo.ApplicationCommandOption{
		// 		{
		// 			Type:        discordgo.ApplicationCommandOptionNumber,
		// 			Name:        "amount",
		// 			Description: "The amount of tokens you want to gamba.",
		// 			MinValue:    &integerOptionMinValue,
		// 			Required:    true,
		// 		},
		// 	},
		// },
		// {
		// 	Name:        "economy",
		// 	Description: "Check the overall number of tokens in the economy",
		// },
		// {
		// 	Name: "invite",
		// 	Description: "Get an invite link to add Bot Person to another server.",
		// },
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
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

				// Logging the interaction to the log file
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, option.StringValue())

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
				persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *i.Interaction.Member.User)

				// Logging outgoing bot response
				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msg,
				})
				if err != nil {

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Something went wrong.")

					// Not 100% sure this is the approach I want to take with handling errors from the API
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong.",
					})
					return
				}
			}
		},
		"my-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< USER_GET_STATS >")

			// Getting user stat data
			msg := persistance.SlashGetUserStats(*i.Interaction.Member.User)

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"bot-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_STATS >")

			// Getting user stat data
			msg := persistance.SlashGetBotStats(s)

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_ABOUT >")

			// Getting user stat data
			msg := "Bot Person started off as a project by AltarCrystal and is now being maintained by Nex. You can see Bot Person's source code at: https://github.com/nexfortisme/bot-person"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"donations": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_DONATIONS >")

			// Getting user stat data
			msg := "Thanks PsychoPhyr for $20 to keep the lights on for Bot Person!"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_HELP >")

			// Getting user stat data
			msg := "A picture is worth 1000 words"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var msg string

			if !persistance.UserHasTokens(i.Interaction.Member.User.ID) {

				persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

				// Logging incoming user request
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< BOT_PERSON_GET_IMAGE >")

				// Getting user stat data
				msg := "You don't have enough tokens to generate an image."

				// Logging outgoing bot response
				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: msg,
					},
				})
				return
			}

			// Pulling the propt out of the optionsMap
			if option, ok := optionMap["prompt"]; ok {

				// Generating the response
				placeholder := "Prompt: " + option.StringValue()

				// Logging incoming user request
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "/image "+option.StringValue())

				// Immediately responding in the 3 second window before the interaciton times out
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: placeholder,
					},
				})

				// Going out to make the OpenAI call to get the proper response
				msg = messages.GetDalleResponseSlashCommand(s, option.StringValue(), config.OpenAIKey)

				persistance.UseImageToken(i.Interaction.Member.User.ID)
				persistance.IncrementInteractionTracking(persistance.BPImageRequest, *i.Interaction.Member.User)

				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &msg,
				})
				if err != nil {

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Something went wrong.")

					// Not 100% sure this is the approach I want to take with handling errors from the API
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong.",
					})
					return
				}
			}
		},
		"balance": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			var tokenCount float64
			var balanceResponse string

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if option, ok := optionMap["user"]; ok {

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_BALANCE > "+option.UserValue(s).Username)

				user := option.UserValue(s)
				tokenCount = persistance.GetUserTokenCount(user.ID)
				balanceResponse = user.Username + " has " + fmt.Sprintf("%.2f", tokenCount) + " tokens."
			} else {

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_BALANCE >")

				tokenCount = persistance.GetUserTokenCount(i.Interaction.Member.User.ID)
				balanceResponse = "You have " + fmt.Sprintf("%.2f", tokenCount) + " tokens."
			}

			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, balanceResponse)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: balanceResponse,
				},
			})
		},
		"send": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			senderBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

			var transferrAmount float64

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// Checking to see that the user has the number of tokens needed to send
			if option, ok := optionMap["amount"]; ok {

				transferrAmount = option.FloatValue()

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_SEND_TOKENS > Amount: "+fmt.Sprintf("%f", transferrAmount))

				if senderBalance < transferrAmount {

					persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Oops! You do not have the tokens needed to complete the transaction.")

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Oops! You do not have the tokens needed to complete the transaction.",
						},
					})
					return
				}
			}

			if option, ok := optionMap["recepient"]; ok {
				recepient := option.UserValue(s)
				sendResponse := persistance.TransferrImageTokens(transferrAmount, i.Interaction.Member.User.ID, recepient.ID)

				newBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_SEND_TOKENS > Amount: "+fmt.Sprintf("%f", transferrAmount)+" Recepient: "+recepient.Username)

				if sendResponse {

					// TODO - Switch to use BPSystemInteraction
					persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Tokens were successfully sent. Your new balance is: "+fmt.Sprint(newBalance))

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Tokens were successfully sent. Your new balance is: " + fmt.Sprint(newBalance),
						},
					})
					return
				} else {

					persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Oops! Something went wrong. Tokens were not sent.")

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Oops! Something went wrong. Tokens were not sent.",
						},
					})
					return
				}
			}
		},
		"bonus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_BONUS >")

			reward, err := persistance.GetUserReward(i.Interaction.Member.User.ID)
			var msg string

			if err != nil {
				msg = err.Error()
			} else {
				// Getting user stat data
				msg = fmt.Sprintf("Congrats! You are awarded %.2f tokens", reward)
			}

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, msg)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})

			// Cleaning up the bonus message if the user is on cooldown
			if err != nil {
				time.Sleep(time.Second * 15)
				s.InteractionResponseDelete(i.Interaction)
			}

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
	flag.BoolVar(&disableCmdReg, "disableCmdReg", false, "Flag for disabling registering of commands on startup")
	flag.Parse()

	readConfig()
	persistance.InitBotStatistics()

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

	if removeCommands {
		removeRegisteredSlashCommands(discordSession)
	}

	if !disableCmdReg {
		registerSlashCommands(discordSession)
	}

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

	// if removeCommands {
	// 	removeRegisteredSlashCommands(discord)
	// }

	persistance.ShutDown()
	_ = discord.Close()
}
