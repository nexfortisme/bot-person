package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Invite(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var retString string
	var url string

	app, err := s.Application("@me")

	if err != nil {
		retString = "Error getting application info"
	}

	logging.LogEvent(eventType.COMMAND_INVITE, i.Interaction.Member.User.ID, "User has requested an invite link", i.Interaction.GuildID)

	url = discordgo.EndpointOAuth2 + "authorize?client_id=" + app.ID + "&permissions=517547084864&scope=bot"
	retString = "Invite me to your server: " + url

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: retString,
		},
	})

}
