package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

type Broken struct{}

func (b *Broken) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "broken",
		Description: "Get more information if something about bot person is broken",
	}
}

func (b *Broken) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logging.LogEvent(eventType.COMMAND_BROKEN, i.Interaction.Member.User.ID, "Broken command used", i.Interaction.GuildID)

	// Getting user stat data
	brokenMessage := "If you have something that is broken about Bot Person, you can create an issue describing what you found here: https://github.com/nexfortisme/bot-person/issues/new or you can join the Bot Person discord and let us know there: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: brokenMessage,
		},
	})
}

func (b *Broken) HelpString() string {
	return "The `/broken` command allows you to get more information if something about Bot Person is broken. You can create an issue on the GitHub repository or join the Bot Person discord and let us know there."
}

func (b *Broken) CommandCost() int {
	return 0
}