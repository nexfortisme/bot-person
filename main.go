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
	"github.com/r3labs/diff/v3"
)

type Config struct {
	OpenAIKey       string `json:"OpenAIKey"`
	DiscordToken    string `json:"DiscordToken"`
	DevDiscordToken string `json:"DevDiscordToken"`
}

var (
	config  Config
	devMode bool

	removeCommands   bool
	removeOnStartup  bool
	removeOnShutdown bool

	disableLogging  bool
	disableTracking bool
	disableCmdReg   bool

	fsInterrupt bool

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
		// TODO - Handle logging of the incoming request by the user
		"bot": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var botResponseString string

			// Pulling the propt out of the optionsMap
			if option, ok := optionMap["prompt"]; ok {

				// Logging the interaction to the log file
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, option.StringValue())

				// Generating the response
				placeholderBotResponse := "Thinking about: " + option.StringValue()

				// Immediately responding in the 3 second window before the interaciton times out
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: placeholderBotResponse,
					},
				})

				// Going out to make the OpenAI call to get the proper response
				botResponseString = messages.ParseSlashCommand(s, option.StringValue(), config.OpenAIKey)

				// Incrementint interaciton counter
				persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *i.Interaction.Member.User)

				// Logging outgoing bot response
				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, botResponseString)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &botResponseString,
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
		"bot-gpt": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var botResponseString string

			// Pulling the propt out of the optionsMap
			if option, ok := optionMap["prompt"]; ok {

				// Logging the interaction to the log file
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, option.StringValue())

				// Generating the response
				placeholderBotResponse := "Thinking about: " + option.StringValue()

				// Immediately responding in the 3 second window before the interaciton times out
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: placeholderBotResponse,
					},
				})

				// Going out to make the OpenAI call to get the proper response
				botResponseString = messages.ParseGPTSlashCommand(s, option.StringValue(), config.OpenAIKey)

				// Incrementint interaciton counter
				persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *i.Interaction.Member.User)

				// Logging outgoing bot response
				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, botResponseString)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &botResponseString,
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
			userStatisticsString := persistance.SlashGetUserStats(*i.Interaction.Member.User)

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, userStatisticsString)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: userStatisticsString,
				},
			})
		},
		"bot-stats": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_STATS >")

			// Getting user stat data
			botStatisticsString := persistance.SlashGetBotStats(s)

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, botStatisticsString)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: botStatisticsString,
				},
			})
		},
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_ABOUT >")

			// Getting user stat data
			aboutMessage := "Bot Person started off as a project by AltarCrystal and is now being maintained by Nex. You can see Bot Person's source code at: https://github.com/nexfortisme/bot-person"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, aboutMessage)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: aboutMessage,
				},
			})
		},
		"donations": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_DONATIONS >")

			// Getting user stat data
			donationMessageString := "Thanks PsychoPhyr for $20 to keep the lights on for Bot Person!"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, donationMessageString)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: donationMessageString,
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_HELP >")

			// Getting user stat data
			helpString := "A picture is worth 1000 words"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, helpString)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: helpString,
				},
			})
		},
		"image": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Access options in the order provided by the user.
			userImageOptions := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			userImageOptionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(userImageOptions))
			for _, opt := range userImageOptions {
				userImageOptionMap[opt.Name] = opt
			}

			var imageReturnString string

			if !persistance.UserHasTokens(i.Interaction.Member.User.ID) {

				persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

				// Logging incoming user request
				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< BOT_PERSON_GET_IMAGE >")

				// Getting user stat data
				imageReturnString := "You don't have enough tokens to generate an image."

				// Logging outgoing bot response
				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, imageReturnString)

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: imageReturnString,
					},
				})
				return
			}

			// Pulling the propt out of the optionsMap
			if option, ok := userImageOptionMap["prompt"]; ok {

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
				imageReturnString = messages.GetDalleResponseSlashCommand(s, option.StringValue(), config.OpenAIKey)

				persistance.UseImageToken(i.Interaction.Member.User.ID)
				persistance.IncrementInteractionTracking(persistance.BPImageRequest, *i.Interaction.Member.User)

				logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, imageReturnString)

				// Updating the initial message with the response from the OpenAI API
				_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &imageReturnString,
				})

				if err != nil {

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Something went wrong.")

					// Not 100% sure this is the approach I want to take with handling errors from the API
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong.",
					})
					return
				} else {
					// TODO - Fix. This doesn't work
					util.CleanUpImages(s, i)
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

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_SEND_TOKENS > Amount: "+fmt.Sprintf("%.2f", transferrAmount))

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

				if i.Interaction.Member.User.ID == recepient.ID {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "You cannot send tokens to yourself.",
						},
					})
					return
				}

				sendResponse := persistance.TransferrImageTokens(transferrAmount, i.Interaction.Member.User.ID, recepient.ID)

				newBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_SEND_TOKENS > Amount: "+fmt.Sprintf("%f", transferrAmount)+" Recepient: "+recepient.Username)

				if sendResponse {

					// TODO - Switch to use BPSystemInteraction
					persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Tokens were successfully sent. Your new balance is: "+fmt.Sprintf("%.2f", newBalance))

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

			bonusReward, returnMessage, err := persistance.GetUserReward(i.Interaction.Member.User.ID)
			var bonusReturnMessage string

			if err != nil {
				bonusReturnMessage = err.Error()
			} else {
				if returnMessage != "" {
					bonusReturnMessage = fmt.Sprintf("%s \nCongrats! You are awarded %.2f tokens", returnMessage, bonusReward)
				} else {
					bonusReturnMessage = fmt.Sprintf("Congrats! You are awarded %.2f tokens", bonusReward)
				}
			}

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, bonusReturnMessage)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: bonusReturnMessage,
				},
			})

			// Cleaning up the bonus message if the user is on cooldown
			if err != nil {
				time.Sleep(time.Second * 15)
				s.InteractionResponseDelete(i.Interaction)
			}

		},
		"lootbox": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_OPEN_LOOTBOX >")

			lootboxReward, lootboxSeed, err := persistance.BuyLootbox(i.Interaction.Member.User.ID)
			var lootboxReturnMessage string

			if err != nil {
				lootboxReturnMessage = err.Error()
			} else {

				// TODO - Refactor this so a change in rates doesn't break the command
				if lootboxReward == 2 {
					lootboxReturnMessage = fmt.Sprintf("%s You purchased a lootbox with the seed: %d and it contained %d tokens", util.GetOofResponse(), lootboxSeed, lootboxReward)
				} else if lootboxReward == 5 {
					lootboxReturnMessage = fmt.Sprintf("You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
				} else if lootboxReward == 20 {
					lootboxReturnMessage = fmt.Sprintf("Congrats! You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
				} else if lootboxReward == 100 {
					lootboxReturnMessage = fmt.Sprintf("Woah! You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
				} else if lootboxReward == 500 {
					lootboxReturnMessage = fmt.Sprintf("Stop Hacking. You purchased a lootbox with the seed: %d and it contained %d tokens", lootboxSeed, lootboxReward)
				}

			}

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, lootboxReturnMessage)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: lootboxReturnMessage,
				},
			})

			// Cleaning up the bonus message if the user is on cooldown
			if err != nil {
				time.Sleep(time.Second * 15)
				s.InteractionResponseDelete(i.Interaction)
			}

		},
		"broken": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

			// Logging incoming user request
			logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_GET_BROKEN >")

			// Getting user stat data
			brokenMessage := "If you have something that is broken about Bot Person, you can create an issue describing what you found here: https://github.com/nexfortisme/bot-person/issues/new"

			// Logging outgoing bot response
			logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, brokenMessage)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: brokenMessage,
				},
			})
		},
		"burn": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			var burnAmount float64
			senderBalance := persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

			// Access options in the order provided by the user.
			options := i.ApplicationCommandData().Options

			// Or convert the slice into a map
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			// Checking to see that the user has the number of tokens needed to send
			if option, ok := optionMap["amount"]; ok {

				burnAmount = option.FloatValue()

				logging.LogIncomingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "< SYSTEM_BURN_TOKENS > Amount: "+fmt.Sprintf("%.2f", burnAmount))

				if senderBalance < burnAmount {

					persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

					logging.LogOutgoingUserInteraction(s, i.Interaction.Member.User.Username, i.Interaction.GuildID, "Oops! You do not have the tokens needed to complete the transaction.")

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Oops! You do not have the tokens needed to complete the transaction.",
						},
					})
					return
				} else {
					persistance.RemoveUserTokens(i.Interaction.Member.User.ID, burnAmount)
					senderBalance = persistance.GetUserTokenCount(i.Interaction.Member.User.ID)

					removeTokenResponse := fmt.Sprintf("%.2f tokens removed. New Balance: %.2f", burnAmount, senderBalance)

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: removeTokenResponse,
						},
					})
				}
			}

		},
	}
)

func readConfig() {

	var botPersonConfig []byte
	botPersonConfig, err := os.ReadFile("config.json")

	if err != nil {
		createdConfig = true
		log.Printf("Error reading config. Creating File")
		os.WriteFile("config.json", []byte("{\"DiscordToken\":\"\",\"OpenAIKey\":\"\"}"), 0666)
		botPersonConfig, err = os.ReadFile("config.json")
		util.HandleFatalErrors(err, "Could not read config file: config.json")
	}

	err = json.Unmarshal(botPersonConfig, &config)
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
	persistance.ReadBotStatistics()

	ticker := time.NewTicker(5 * time.Minute)

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
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill, os.Interrupt)

	for {
		select {
		case <-ticker.C:
			saveBotStatistics()
		case <-interrupt:
			fmt.Println("Interrupt received, stopping...")
			ticker.Stop()
			shutDown(discordSession)
			return
		}
	}

}

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
