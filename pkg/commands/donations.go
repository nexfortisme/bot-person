package commands

import (
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func Donations(s *discordgo.Session, i *discordgo.InteractionCreate) {
	
	logging.LogEvent(eventType.COMMAND_DONATIONS, i.Interaction.Member.User.ID, "Donations command used", i.Interaction.GuildID)

	// Getting user stat data
	donationMessageString := "Thanks PsychoPhyr for $20 to keep the lights on for Bot Person!\n If you would like to contribute, you can do so in the Bot Person Discord Server: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: donationMessageString,
		},
	})
}
