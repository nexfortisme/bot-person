package commands

import "github.com/bwmarrin/discordgo"

func Invite(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var retString string
	var url string

	app, err := s.Application("@me")

	if err != nil {
		retString = "Error getting application info"
	}

	url = discordgo.EndpointOAuth2 + "authorize?client_id=" + app.ID + "&permissions=517547084864&scope=bot"
	retString = "Invite me to your server: " + url

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: retString,
		},
	})

}
