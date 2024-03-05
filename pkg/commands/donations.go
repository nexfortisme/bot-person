package commands

import (
	loggingType "main/pkg/logging/enums"
	logging "main/pkg/logging/services"

	"github.com/bwmarrin/discordgo"
)

func Donations(s *discordgo.Session, i *discordgo.InteractionCreate) {
	logging.LogEvent(loggingType.COMMAND_DONATIONS, "User requested donation information", i.Interaction.Member.User.Username, i.Interaction.GuildID, s)

	donationMessageString := "Thanks PsychoPhyr for helping to keep the lights on for Bot Person!\n If you would like to contribute, you can do so in the Bot Person Discord Server: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: donationMessageString,
		},
	})
}
