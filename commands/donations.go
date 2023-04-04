package commands

import (
	"main/persistance"

	"github.com/bwmarrin/discordgo"
)

func Donations(s *discordgo.Session, i *discordgo.InteractionCreate) {
	persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *i.Interaction.Member.User)

	// Getting user stat data
	donationMessageString := "Thanks PsychoPhyr for $20 to keep the lights on for Bot Person!"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: donationMessageString,
		},
	})
}