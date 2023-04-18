package messages

import (
	"main/external"
	"main/persistance"
	"main/util"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	var incomingMessage string

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Checking for prefix
	if !strings.HasPrefix(m.Message.Content, "!") {
		incomingMessage = strings.ToLower(m.Message.Content)
	} else {
		incomingMessage = m.Message.Content
	}

	persistance.APictureIsWorthAThousand(incomingMessage, m)

	// Checking for specific commands: !addTokens, !setGPT3, !setGPT4
	if strings.HasPrefix(incomingMessage, "!addTokens") {

		if !util.UserIsAdmin(m.Author.ID) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		}

		req := strings.Split(incomingMessage, " ")

		if len(req) != 3 {
			return
		}

		tokenCount, _ := strconv.ParseFloat(req[2], 64)
		success := persistance.AddBotPersonTokens(tokenCount, req[1][2:len(req[1])-1])

		if success {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Tokens were successfully added.")
		} else {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Something went wrong. Tokens were not added.")
		}
	} else if strings.HasPrefix(incomingMessage, "!setGPT4") {

		if !util.UserIsAdmin(m.Author.ID) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		}

		external.SetGPT4()
		_, _ = s.ChannelMessageSend(m.ChannelID, "Model set to GPT-4")
	} else if strings.HasPrefix(incomingMessage, "!setGPT3") {

		if !util.UserIsAdmin(m.Author.ID) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		}

		external.SetGPT3()
		_, _ = s.ChannelMessageSend(m.ChannelID, "Model set to GPT-3")

	}

}
