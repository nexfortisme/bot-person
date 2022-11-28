package messages

import (
	"fmt"
	"main/logging"
	"main/messages/external"
	"main/persistance"
	"main/util"
	"math/rand"
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
	badBotResponses := make([]string, 0)
	badBotResponses = append(badBotResponses, "I'm sorry")
	badBotResponses = append(badBotResponses, "It won't happen again")
	badBotResponses = append(badBotResponses, "Eat Shit")
	badBotResponses = append(badBotResponses, "Ok.")
	badBotResponses = append(badBotResponses, "Sure Thing.")
	badBotResponses = append(badBotResponses, "Like you are the most perfect being in existance. Pound sand pal.")
	badBotResponses = append(badBotResponses, "https://youtu.be/4X7q87RDSHI")

	if !strings.HasPrefix(m.Message.Content, "!") {
		incomingMessage = strings.ToLower(m.Message.Content)
	} else {
		incomingMessage = m.Message.Content
	}

	// Looking at messages in the channel and returning WORD_COUNT / 1000 number of tokens
	// ie. A picture is worth 1000 words
	wordCount := len(strings.Fields(incomingMessage))
	tokenValue := fmt.Sprintf("%.2f", (float64(wordCount) / 1000.0))
	tokenAddAmount, _ := strconv.ParseFloat(tokenValue, 64)

	persistance.AddImageTokens(tokenAddAmount, *&m.Author.ID)

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {
		logging.LogIncomingMessage(s, m)
		persistance.IncrementInteractionTracking(persistance.BPBadBotInteraction, *m.Author)
		badBotRetort := badBotResponses[rand.Intn(len(badBotResponses))]
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
	} else if strings.HasPrefix(incomingMessage, ";;lenny"){
		persistance.IncrementInteractionTracking(persistance.BPLennyFaceInteracton, *m.Author)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)");
	}
	// } else if strings.HasPrefix(incomingMessage, "!gamble") {

	// 	// !TODO - FIX ISSUE WITH FLOATING POINT

	// 	persistance.IncrementInteractionTracking(persistance.BPBadBotInteraction, *m.Author)
	// 	s.ChannelMessageSend(m.ChannelID, "Fuck Off.")

	// 	// req := strings.Split(incomingMessage, " ")
	// 	// tokenCount, _ := strconv.ParseFloat(req[1], 64)
	// 	// persistance.RemoveUserTokens(m.Author.ID, tokenCount)

	// 	// rand.Seed(time.Now().UnixNano())
	// 	// num := rand.Intn(101)

	// 	// retStr := "Bot Person Rolled a " + strconv.Itoa(num) + "."

	// 	// if num < 50 {
	// 	// 	retStr += " OOF. You lose the tokens you gambled. :("
	// 	// 	s.ChannelMessageSend(m.ChannelID, retStr)
	// 	// 	persistance.RemoveUserTokens(m.Author.ID, tokenCount)
	// 	// } else if num >= 50 && num < 80 {
	// 	// 	retStr += " Nice Profit! You win 1.1x what you gambled."
	// 	// 	s.ChannelMessageSend(m.ChannelID, retStr)
	// 	// 	rewards := tokenCount * 1.1
	// 	// 	persistance.AddImageTokens(math.Floor(rewards*100)/100, m.Author.ID)
	// 	// } else if num >= 80 && num < 90 {
	// 	// 	retStr += " Good Profit! You win 1.2x what you gambled."
	// 	// 	s.ChannelMessageSend(m.ChannelID, retStr)
	// 	// 	rewards := tokenCount * 1.2
	// 	// 	persistance.AddImageTokens(math.Floor(rewards*100)/100, m.Author.ID)
	// 	// } else if num >= 90 && num <= 99 {
	// 	// 	retStr += " Great Profit! You win 1.4x what you gambled."
	// 	// 	s.ChannelMessageSend(m.ChannelID, retStr)
	// 	// 	rewards := tokenCount * 1.4
	// 	// 	persistance.AddImageTokens(math.Floor(rewards*100)/100, m.Author.ID)
	// 	// } else if num > 99 {
	// 	// 	retStr += " Jackpot! You win 2x what you gambled."
	// 	// 	s.ChannelMessageSend(m.ChannelID, retStr)
	// 	// 	rewards := tokenCount * 2
	// 	// 	persistance.AddImageTokens(math.Floor(rewards*100)/100, m.Author.ID)
	// 	// }

	// }

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
	return strings.HasPrefix(m.Content, "!bot") && m.Content != "!botStats"
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
