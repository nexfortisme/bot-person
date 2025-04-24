package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type Help struct {}

func (h *Help) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
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
	}
}

func (h *Help) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	logging.LogEvent(eventType.COMMAND_HELP, i.Interaction.Member.User.ID, "Help command used", i.Interaction.GuildID)

	var helpOption string
	var helpString string

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["command"]; ok {
		helpOption = option.StringValue()
	}

	// Use a switch statement to determine the help message based on the command
	switch helpOption {
	case "":
		helpString = (&Help{}).HelpString()
	case "bot":
		helpString = (&Bot{}).HelpString()
	case "bot-gpt":
		helpString = (&BotGPT{}).HelpString()
	case "my-stats":
		helpString = (&MyStats{}).HelpString()
	case "bot-stats":
		helpString = (&BotStats{}).HelpString()
	case "about":
		helpString = (&About{}).HelpString()
	case "donations":
		helpString = (&Donations{}).HelpString()
	case "images":
		helpString = (&Image{}).HelpString()
	case "store":
		helpString = "DEPRECATED"
	case "balance":
		helpString = (&Balance{}).HelpString()
	case "send":
		helpString = (&Send{}).HelpString()
	case "bonus":
		helpString = (&Bonus{}).HelpString()
	case "loot-box":
		helpString = (&Lootbox{}).HelpString()
	case "invite":
		helpString = (&Invite{}).HelpString()
	case "burn":
		helpString = (&Burn{}).HelpString()
	case "stocks":
		helpString = "DEPRECATED"
	case "portfolio":
		helpString = "DEPRECATED"
	case "broken":
		helpString = (&Broken{}).HelpString()
	case "save-streak":
		helpString = "DEPRECATED"
	case "hsr-code":
		helpString = (&HSRCode{}).HelpString()
	case "search":
		helpString = (&Search{}).HelpString()
	default:
		helpString = "Command not found. Use `/help` for a list of available commands."
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpString,
		},
	})
}

func (h *Help) HelpString() string {
	return "To get help with a specific command use `/help [command]` where `[command]` is the slash command you want help with. If that still doesn't help, feel free to join the Bot Person discord server and ask there and someone will try their best to sort out what ever issue you may have. https://discord.gg/MtEG5zMtUR"
}

func (h *Help) CommandCost() int {
	return 0
}
