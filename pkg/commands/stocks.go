package commands

import (
	"github.com/bwmarrin/discordgo"
)

func Stocks(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Stocks command is deprecated. Not going to be re-implemented.",
		},
	})
}
