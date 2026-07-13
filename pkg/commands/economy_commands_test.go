package commands

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestTokenAmountOptionsArePositiveIntegers(t *testing.T) {
	commands := []*discordgo.ApplicationCommand{
		(&Send{}).ApplicationCommand(),
		(&Burn{}).ApplicationCommand(),
	}

	for _, command := range commands {
		t.Run(command.Name, func(t *testing.T) {
			var amountOption *discordgo.ApplicationCommandOption
			for _, option := range command.Options {
				if option.Name == "amount" {
					amountOption = option
					break
				}
			}
			if amountOption == nil {
				t.Fatal("amount option not found")
			}
			if amountOption.Type != discordgo.ApplicationCommandOptionInteger {
				t.Fatalf("amount option type = %v, want integer", amountOption.Type)
			}
			if amountOption.MinValue == nil || *amountOption.MinValue < 1 {
				t.Fatalf("amount minimum = %v, want at least 1", amountOption.MinValue)
			}
		})
	}
}
