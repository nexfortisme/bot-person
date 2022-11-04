package messages

import (
	"log"
	"main/logging"
	"main/messages/external"
	"main/util"
	"math/rand"
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

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(incomingMessage, "bad bot") {
		logging.LogIncomingMessage(s, m)

		logging.IncrementTracker(2, m.Author.ID, m.Author.Username)
		badBotRetort := badBotResponses[rand.Intn(len(badBotResponses))]
		// TODO - Here Too
		log.Println("Bot Person > " + badBotRetort)
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "good bot") {
		logging.LogIncomingMessage(s, m)

		logging.IncrementTracker(1, m.Author.ID, m.Author.Username)
		log.Println("Bot Person > Thank You!")
		_, err := s.ChannelMessageSend(m.ChannelID, "Thank You!")
		util.HandleErrors(err)
	} else if strings.HasPrefix(incomingMessage, "!botStats") {
		logging.LogIncomingMessage(s, m)

		logging.GetBotStats(s, m)
		logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	} else if strings.HasPrefix(incomingMessage, "!myStats") {
		logging.LogIncomingMessage(s, m)

		logging.GetUserStats(s, m)
		logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	} else if strings.HasPrefix(incomingMessage, "!image") {

		logging.LogIncomingMessage(s, m);
		logging.IncrementTracker(3, m.Author.ID, m.Author.Username);

		req := strings.SplitAfterN(incomingMessage, " ", 2);
		resp := external.GetDalleResponse(req[1], openAIKey);
		_, err := s.ChannelMessageSend(m.ChannelID, resp)
		util.HandleErrors(err);
	}

	// Commands to add
	// about - list who made it and maybe a link to the git repo
	// invite - Generates an invite link to be able to invite the bot to differnet servers
	// stopTracking - Allows uers to opt out of data collection

	// Only process messages that mention the bot
	id := s.State.User.ID
	if !mentionsBot(m.Mentions, id) && !mentionsKeyphrase(m) {
		return
	}

	msg := util.ReplaceIDsWithNames(m, s)

	logging.LogIncomingMessage(s, m)

	logging.IncrementTracker(0, m.Author.ID, m.Author.Username)
	respTxt := external.GetOpenAIResponse(msg, openAIKey)

	// TODO - Here as well
	log.Printf("Bot Person > %s \n", respTxt)
	if mentionsKeyphrase(m) {
		s.ChannelMessageSend(m.ChannelID, "!bot is deprecated. Please at the bot or use /bot for further interactions")
	}
	_, err := s.ChannelMessageSend(m.ChannelID, respTxt)
	util.HandleErrors(err)

}

// TODO - Make the response that is being logged by the bot include the bot user's actual username instead of "Bot Person"
func ParseSlashCommand(s *discordgo.Session, prompt string, openAIKey string) string {
	respTxt := external.GetOpenAIResponse(prompt, openAIKey)
	respTxt = "Request: " + prompt + " " + respTxt
	log.Printf("Bot Person > %s \n", respTxt)
	return respTxt
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
