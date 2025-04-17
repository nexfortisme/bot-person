package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Help(s *discordgo.Session, i *discordgo.InteractionCreate) {

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
		helpString = "To get help with a specific command use `/help [command]` where `[command]` is the slash command you want help with. If that still doesn't help, feel free to join the Bot Person discord server and ask there and someone will try their best to sort out what ever issue you may have. https://discord.gg/MtEG5zMtUR"
	case "bot":
		helpString = "The `/bot` command allows you to prompt OpenAI's Divinci chat model. You can ask it whatever as part of the `prompt` and once it generates a response, it will update the message with what came back."
	case "bot-gpt":
		helpString = "The `/bot-gpt` command allows you to prompt OpenAI's GPT-3 or GPT-4 chat model. You can ask it whatever as part of the `prompt` and once it generates a response, it will update the message with what came back. This is slower than the `/bot` command due to the chat model being more complex."
	case "my-stats":
		helpString = "The `/my-stats` command allows you to see your current stats. This includes the number of interactions you have had with the bot, the number of images requested from the `/image` command, your current `/bonus` streak, the number of save streak tokens you have, and what stocks you currently own, if any."
	case "about":
		helpString = "The `/about` command gives a small backstory about Bot Person and links out to the GitHub repository and the Bot Person discord server."
	case "donations":
		helpString = "The `/donations` command gives credit to those who have donated to keeping the lights on for Bot Person and gives further information for those who wish to contribute."
	case "images":
		helpString = "The `/images` command allows you to request an image from OpenAI's Dall-E API at the cost of 1 Bot Person token per image. The image returned is based on what you give in the `prompt` option."
	case "store":
		helpString = "The `/store` command allows you to purchase items/goods to use directly with Bot Person or with other users. You can specifiy which item to purchase with the `item` option. Currently the items that can be purchased are: \nSave Streak Token (50 Tokens)"
	case "balance":
		helpString = "The `/balance` command allows you to see how many Bot Person tokens you currently have. You can also see the balance of others in the server with the `user` option. If you don't specify a user, it will default to you."
	case "send":
		helpString = "The `/send` command allows you to send Bot Person tokens to other users in the server. You can specify the user to send to with the `user` option and the amount to send with the `amount` option."
	case "bonus":
		helpString = "The `/bonus` command allows you to claim your daily bonus tokens. You can only claim this once in a 24 hour period. There are greater rewards for keeping a streak alive. If you miss a day, you will be offered instructions to save your streak through the `/save-streak` command."
	case "loot-box":
		helpString = "The `/loot-box` command allows you to open a loot box for 2.5 Bot Person Tokens and receive between 2 and 500 Bot Person tokens as a reward."
	case "invite":
		helpString = "The `/invite` command generates an invite link with the proper permissions to invite Bot Person to your server."
	case "burn ":
		helpString = "The `/burn` command allows you to remove a specified number of Bot Person tokens from your balance with the `amount` option. This is irreversible and you will not be able to get those tokens back."
	case "stocks":
		helpString = "The `/stock` command allows you to buy and sell stocks with Bot Person tokens at the cost of 1 USD = 1 Bot Person Token. These stocks are entirely fictional and their only purpose is to offer different means of getting Bot Person Tokens beyond that of the `/bonus` command.\nTo purchase stocks, you specity that you want to `Buy` in the `action` option, specify the stock ticker with the `stock` and the quantity with the `quantity` option. You can purchase as few as .1 shares and as many as your balance would allow. To sell, its similar to buying but you specify `Sell` in the `action` option and follow the same steps as buying for the stock and quantity. You cannot sell a stock you don't have."
	case "portfolio":
		helpString = "The `/portfolio` command allows you to see your current portfolio. This includes the stocks you currently own, the quantity of each stock. This information can also be seen in `/my-stats`."
	case "broken":
		helpString = "The `/broken` command gives information on how to report issues with Bot Person."
	case "save-streak":
		helpString = "The `/save-streak` command is only available once you have failed to collect your `/bonus` within a 24 hour period. You have the options of `use` or `buy`. The `use` option will use any Save Streak Tokens you currently possess and continue your streak there. The `buy` option will also allow you to keep your streak but due to the dire nature of continuing your streak, it will cost ***HALF*** of your current balance."
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
