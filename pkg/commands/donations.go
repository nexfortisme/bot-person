package commands

import (

	"github.com/bwmarrin/discordgo"
)

type Donations struct{}

func (d *Donations) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "donations",
		Description: "List of the people who contributed to Bot Person's on-going service.",
	}
}

func (d *Donations) Execute(s *discordgo.Session, i *discordgo.InteractionCreate) {

	// Getting user stat data
	donationMessageString := "Thanks PsychoPhyr for $20 to keep the lights on for Bot Person!\n If you would like to contribute, you can do so in the Bot Person Discord Server: https://discord.gg/MtEG5zMtUR"

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: donationMessageString,
		},
	})
}

func (d *Donations) HelpString() string {
	return "The `/donations` command gives credit to those who have donated to keeping the lights on for Bot Person and gives further information for those who wish to contribute."
}

func (d *Donations) CommandCost() int {
	return 0
}
