package commands

import (
	"main/pkg/logging"

	eventType "main/pkg/logging/enums"


	"github.com/bwmarrin/discordgo"
)

type About struct {}

func (a *About) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "about",
		Description: "Get information about Bot Person.",
	}
}

func (a *About) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.COMMAND_ABOUT, i.Interaction.Member.User.ID, "About command used", i.Interaction.GuildID)

	// Getting user stat data
	aboutMessage := "Bot Person started off as a project by AltarCrystal and is now being maintained by Nex. You can see Bot Person's source code at: https://github.com/nexfortisme/bot-person or learn more at the Bot Person discord: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: aboutMessage,
		},
	})
}

func (a *About) HelpString() string {
	return "The `/about` command gives a small backstory about Bot Person and links out to the GitHub repository and the Bot Person discord server."
}

func (a *About) CommandCost() int {
	return 0
}
