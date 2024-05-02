package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"
	state "main/pkg/state/services"

	"github.com/bwmarrin/discordgo"
)

func Leave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.TTS_JOIN, i.Interaction.Member.User.ID, fmt.Sprintf("Bot Joined Channel: %s", i.Interaction.ChannelID), i.Interaction.GuildID)
	
	external.Leave(state.GetConnection(i.Interaction.ChannelID));
	state.RemoveConnection(i.Interaction.ChannelID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Bot Left Channel",
		},
	})
}
