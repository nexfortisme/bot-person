package messages

import (
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
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

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {

		persistance.IncrementInteractionTracking(persistance.BPBadBotInteraction, *m.Author)

		badBotRetort := util.GetBadBotResponse()

		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "good bot") {
		persistance.IncrementInteractionTracking(persistance.BPGoodBotInteraction, *m.Author)

		goodBotRetort := util.GetGoodBotResponse()

		_, err := s.ChannelMessageSend(m.ChannelID, goodBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "!addTokens") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		if !util.UserIsAdmin(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			req := strings.Split(incomingMessage, " ")

			if len(req) != 3 {
				return
			}

			tokenCount, _ := strconv.ParseFloat(req[2], 64)
			success := persistance.AddBotPersonTokens(tokenCount, req[1][2:len(req[1])-1])
			if success {
				s.ChannelMessageSend(m.ChannelID, "Tokens were successfully added.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Something went wrong. Tokens were not added.")
			}
		}
	} else if strings.HasPrefix(incomingMessage, "!setGPT4") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		if !util.UserIsAdmin(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			external.SetGPT4()
			s.ChannelMessageSend(m.ChannelID, "Model set to GPT-4")
		}
	} else if strings.HasPrefix(incomingMessage, "!setGPT3") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		if !util.UserIsAdmin(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			external.SetGPT3()
			s.ChannelMessageSend(m.ChannelID, "Model set to GPT-3")
		}
	} else if strings.HasPrefix(incomingMessage, "!setGPT4New") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		if !util.UserIsAdmin(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			external.SetGPT4New()
			s.ChannelMessageSend(m.ChannelID, "Model set to GPT-4")
		}
	} else if strings.HasPrefix(incomingMessage, ";;lenny") {
		persistance.IncrementInteractionTracking(persistance.BPLennyFaceInteracton, *m.Author)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
	}

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	msg := util.ReplaceIDsWithNames(m, s)

	persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *m.Author)
	respTxt := external.GetOpenAIResponse(msg)

	if mentionsKeyphrase(m) {
		s.ChannelMessageSend(m.ChannelID, "!bot is deprecated. Please at the bot or use /bot for further interactions")
	}
	_, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	util.HandleErrors(err)
}

func mentionsKeyphrase(m *discordgo.MessageCreate) bool {
	return strings.HasPrefix(m.Content, "!bot")
}

// Determine if the bot's ID is in the list of users mentioned
func mentionsBot(mentions []*discordgo.User, id string) bool {
	for _, user := range mentions {
		if user.ID == id {
			return true
		}
	}
	return false
}
