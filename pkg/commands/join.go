package commands

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"
	state "main/pkg/state/services"

	"github.com/bwmarrin/discordgo"
)

func Join(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.TTS_JOIN, i.Interaction.Member.User.ID, fmt.Sprintf("Bot Joined Channel: %s", i.Interaction.ChannelID), i.Interaction.GuildID)
	
	state.SetConnection(i.Interaction.ChannelID, external.Join(s, i.Interaction))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Bot Joined Channel",
		},
	})
}
