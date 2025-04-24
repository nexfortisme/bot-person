package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	ApplicationCommand() *discordgo.ApplicationCommand
	Execute(s *discordgo.Session, i *discordgo.InteractionCreate)
	HelpString() string

	CommandCost() int
}
