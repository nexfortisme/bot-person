package commands

import (
	"main/pkg/logging"

	eventType "main/pkg/logging/enums"


	"github.com/bwmarrin/discordgo"
)

func About(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
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
