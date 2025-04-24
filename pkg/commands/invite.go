package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type Invite struct{}

func (in *Invite) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "invite",
		Description: "Get an invite link to invite Bot Person to your server.",
	}
}

func (in *Invite) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var retString string
	var url string

	app, err := s.Application("@me")

	if err != nil {
		retString = "Error getting application info"
	}

	logging.LogEvent(eventType.COMMAND_INVITE, i.Interaction.Member.User.ID, "User has requested an invite link", i.Interaction.GuildID)

	// TODO - Is there a way to calculate the permissions on the fly?
	url = discordgo.EndpointOAuth2 + "authorize?client_id=" + app.ID + "&permissions=517547084864&scope=bot"
	retString = "Invite me to your server: " + url

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: retString,
		},
	})

}

func (in *Invite) HelpString() string {
	return "The `/invite` command generates an invite link with the proper permissions to invite Bot Person to your server."
}

func (in *Invite) CommandCost() int {
	return 0
}
