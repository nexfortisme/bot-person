package commands

import "github.com/bwmarrin/discordgo"

func Leaderboard(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var leaderboardOption string
	var returnString string

	// Access options in the order provided by the user.
	options := i.ApplicationCommandData().Options

	// Or convert the slice into a map
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Checking to see that the user has the number of tokens needed to send
	if option, ok := optionMap["action"]; ok {
		leaderboardOption = option.StringValue()
	}

	// Don't have to worry about the else case since the slash command will force the user to pick one of two options
	if leaderboardOption == "tokens" {
		
	} else if leaderboardOption == "streaks" {

	}

}
