package commands

import (
	// "main/pkg/persistance"

	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
)

func Broken(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logging.LogEvent(loggingType.COMMAND_BROKEN, "User reported a broken feature", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

	brokenMessage := "If you have something that is broken about Bot Person, you can create an issue describing what you found here: https://github.com/nexfortisme/bot-person/issues/new or you can join the Bot Person discord and let us know there: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: brokenMessage,
		},
	})
}
