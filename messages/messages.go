package messages

import (
	"main/logging"
	"main/messages/external"
	"main/persistance"
	"main/util"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate, openAIKey string) {

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	// TODO - Add this to the config file
	var incomingMessage string

	if !strings.HasPrefix(m.Message.Content, "!") {
		incomingMessage = strings.ToLower(m.Message.Content)
	} else {
		incomingMessage = m.Message.Content
	}

	persistance.APictureIsWorthAThousand(incomingMessage, m)

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {

		logging.LogIncomingMessage(s, m)
		persistance.IncrementInteractionTracking(persistance.BPBadBotInteraction, *m.Author)

		badBotRetort := util.GetBadBotResponse()

		logging.LogOutgoingUserInteraction(s, m.Author.Username, m.GuildID, badBotRetort)

		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)

	} else if strings.HasPrefix(incomingMessage, "good bot") {
		logging.LogIncomingMessage(s, m)
		persistance.IncrementInteractionTracking(persistance.BPGoodBotInteraction, *m.Author)

		logging.LogOutgoingUserInteraction(s, m.Author.Username, m.GuildID, "Thank You!")

		_, err := s.ChannelMessageSend(m.ChannelID, "Thank You!")
		util.HandleErrors(err)

	} else if strings.HasPrefix(incomingMessage, "!addTokens") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		// TODO - Change this to pull from the config instead of being a hardcoded value
		if m.Author.ID != "92699061911580672" {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			req := strings.Split(incomingMessage, " ")
			tokenCount, _ := strconv.ParseFloat(req[2], 64)
			success := persistance.AddImageTokens(tokenCount, req[1][2:len(req[1])-1])
			if success {
				s.ChannelMessageSend(m.ChannelID, "Tokens were successfully added.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Something went wrong. Tokens were not added.")
			}
		}
	} else if strings.HasPrefix(incomingMessage, "!merryChristmas") {

		// TODO - Switch to use BPSystemInteraction
		persistance.IncrementInteractionTracking(persistance.BPBasicInteraction, *m.Author)

		// TODO - Change this to pull from the config instead of being a hardcoded value
		if m.Author.ID != "92699061911580672" {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			for _, element := range persistance.GetUserStats() {
				element.UserStats.ImageTokens += 25;
				if !persistance.UpdateUser(element) {
					return
				}
			}
			s.ChannelMessageSend(m.ChannelID, "Merry Christmas! 25 Tokens Added To All Users")
		}
	} else if strings.HasPrefix(incomingMessage, ";;lenny") {
		persistance.IncrementInteractionTracking(persistance.BPLennyFaceInteracton, *m.Author)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
	}

	// ! Add Help Command

	// Commands to add
	// invite - Generates an invite link to be able to invite the bot to differnet servers
	// stopTracking - Allows uers to opt out of data collection

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	msg := util.ReplaceIDsWithNames(m, s)

	logging.LogIncomingMessage(s, m)

	persistance.IncrementInteractionTracking(persistance.BPChatInteraction, *m.Author)
	respTxt := external.GetOpenAIResponse(msg, openAIKey)

	logging.LogOutgoingUserInteraction(s, m.Author.Username, m.GuildID, respTxt)

	if mentionsKeyphrase(m) {
		s.ChannelMessageSend(m.ChannelID, "!bot is deprecated. Please at the bot or use /bot for further interactions")
	}
	_, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	util.HandleErrors(err)

}

func ParseSlashCommand(s *discordgo.Session, prompt string, openAIKey string) string {
	respTxt := external.GetOpenAIResponse(prompt, openAIKey)
	respTxt = "Request: " + prompt + " " + respTxt
	return respTxt
}

// TODO - Rename, I don't like this
func GetDalleResponseSlashCommand(s *discordgo.Session, prompt string, openAIKey string) string {
	dalleResponse, err := external.GetDalleResponse(prompt, openAIKey)

	if err != nil {
		return dalleResponse
	}

	dalleResponse = "Prompt: " + "[" + prompt + "](" + dalleResponse + ")"
	return dalleResponse
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
