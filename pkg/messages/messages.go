package messages

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"strconv"
	"strings"

	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

var connections = make(map[string]*discordgo.VoiceConnection)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	var incomingMessage string
	var isReply bool = false
	var originalMessage string

	// This means the current message is a reply to another message
	if m.Message.ReferencedMessage != nil {
		isReply = true
		originalMessage = m.Message.ReferencedMessage.Content
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

		logging.LogEvent(eventType.USER_BAD_BOT, m.Author.ID, "Bad bot command used", m.GuildID)

		badBotRetort := util.GetBadBotResponse()

		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "good bot") {

		logging.LogEvent(eventType.USER_GOOD_BOT, m.Author.ID, "Good bot command used", m.GuildID)

		goodBotRetort := util.GetGoodBotResponse()

		_, err := s.ChannelMessageSend(m.ChannelID, goodBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "!addTokens") {

		if !util.UserIsAdmin(m.Author.ID) {

			logging.LogEvent(eventType.ECONOMY_CREATE_TOKENS, m.Author.ID, "NOT ENOUGH PERMISSIONS", m.GuildID)

			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {

			logging.LogEvent(eventType.ECONOMY_CREATE_TOKENS, m.Author.ID, "Add tokens command used", m.GuildID)

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
	} else if strings.HasPrefix(incomingMessage, ";;lenny") {

		logging.LogEvent(eventType.LENNY, m.Author.ID, "Lenny command used", m.GuildID)

		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
	}

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	if isReply && mentionsBot(m.Mentions, id) {
		// Remove any @ mentions from the messages
		cleanedIncomingMessage := strings.ReplaceAll(incomingMessage, "<@"+id+"> ", "")
		cleanedOriginalMessage := strings.ReplaceAll(originalMessage, "<@"+id+">", "")

		perplexityResponse := external.GetPerplexityResponse(cleanedOriginalMessage, cleanedIncomingMessage)
		response := perplexityResponse.Choices[0].Message.Content

		if perplexityResponse.Citations != nil {
			for index, citation := range perplexityResponse.Citations {
				replaceString := fmt.Sprintf("[%d]", index)
				replacementString := fmt.Sprintf("[[%d]](%s)", index, citation)
				response = strings.Replace(response, replaceString, replacementString, 1)
			}
		}
		s.ChannelMessageSendReply(m.ChannelID, response, m.Message.Reference())
	}
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
