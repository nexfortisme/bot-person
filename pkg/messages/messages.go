package messages

import (
	"encoding/json"
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"strconv"
	"strings"
	"time"

	"main/pkg/commands"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// -------------------- DEBUGGING --------------------
	// currentTime := time.Now().Format("2006-01-02-15-04-05")

	// data, _ := json.MarshalIndent(m, "", "  ")
	// util.SaveResponseToFile(data, fmt.Sprintf("message-%s.txt", currentTime))
	// -------------------- DEBUGGING --------------------

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	var isReply bool = false
	var replyType string

	// This means the current message is a reply to another message
	if m.Message.ReferencedMessage != nil {
		isReply = true

		if( m.ReferencedMessage.Interaction == nil ) {
			replyType = "message"
			return
		}

		switch m.ReferencedMessage.Interaction.Name {
		case "bot":
			replyType = "bot"
		case "bot_gpt":
			replyType = "bot_gpt"
		case "image":
			replyType = "image"
		default:
			replyType = "message"
		}
	}

	// persistance.APictureIsWorthAThousand(m.Message.Content, m)

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(m.Message.Content, "bad bot") {
		logging.LogEvent(eventType.USER_BAD_BOT, m.Author.ID, "Bad bot command used", m.GuildID)
		badBotRetort := util.GetBadBotResponse()
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(m.Message.Content, "good bot") {
		logging.LogEvent(eventType.USER_GOOD_BOT, m.Author.ID, "Good bot command used", m.GuildID)
		goodBotRetort := util.GetGoodBotResponse()
		_, err := s.ChannelMessageSend(m.ChannelID, goodBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(m.Message.Content, "!addTokens") {
		if !util.UserIsAdmin(m.Author.ID) {
			logging.LogEvent(eventType.ECONOMY_CREATE_TOKENS, m.Author.ID, "NOT ENOUGH PERMISSIONS", m.GuildID)
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		} else {
			logging.LogEvent(eventType.ECONOMY_CREATE_TOKENS, m.Author.ID, "Add tokens command used", m.GuildID)
			req := strings.Split(m.Message.Content, " ")
			if len(req) != 3 {
				return
			}

			tokenCount, _ := strconv.ParseInt(req[2], 10, 64)
			success := persistance.AddBotPersonTokens(int(tokenCount), req[1][2:len(req[1])-1])
			if success {
				s.ChannelMessageSend(m.ChannelID, "Tokens were successfully added.")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Something went wrong. Tokens were not added.")
			}
		}
	} else if strings.HasPrefix(m.Message.Content, ";;lenny") {
		logging.LogEvent(eventType.LENNY, m.Author.ID, "Lenny command used", m.GuildID)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
	}

	// Reply Handling
	if isReply && replyType == "message" {

		// If it is a message response and doesn't mention the bot, we out
		if !mentionsBot(m.Mentions, s.State.User.ID) {
			return
		}

		cleanedIncomingMessage := strings.ReplaceAll(m.Message.Content, "<@"+s.State.User.ID+"> ", "")
		cleanedOriginalMessage := strings.ReplaceAll(m.ReferencedMessage.Content, "<@"+s.State.User.ID+">", "")

		perplexityResponse := external.GetPerplexityResponse(cleanedOriginalMessage, cleanedIncomingMessage)

		if len(perplexityResponse.Choices) == 0 {
			s.ChannelMessageSendReply(m.ChannelID, "Error getting response from Perplexity.", m.Message.Reference())

			currentTime := time.Now().Format("2006-01-02-15-04-05")

			data, _ := json.MarshalIndent(perplexityResponse, "", "  ")
			util.SaveResponseToFile(data, fmt.Sprintf("perplexity-error-response-%s.txt", currentTime))

			return
		}

		response := perplexityResponse.Choices[0].Message.Content

		if perplexityResponse.Citations != nil {
			for index, citation := range perplexityResponse.Citations {
				replaceString := fmt.Sprintf("[%d]", index)
				replacementString := fmt.Sprintf("[[%d]](%s)", index, citation)
				response = strings.Replace(response, replaceString, replacementString, 1)
			}
		}
		s.ChannelMessageSendReply(m.ChannelID, response, m.Message.Reference())

		return
	} else if isReply && replyType == "bot" {
		return
	} else if isReply && replyType == "bot_gpt" {
		return
	} else if isReply && replyType == "image" {

		originalImageId := m.ReferencedMessage.Attachments[0].Filename
		originalImageId = strings.Split(originalImageId, ".")[0] // Removes extension. File name is <openai_message_id>.jpg

		followUpImagePrompt := util.EscapeQuotes(m.Message.Content)

		user, err := persistance.GetUser(m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong. Follow up image was not generated.")
			return
		}

		if user.ImageTokens < (&commands.Image{}).CommandCost() {
			s.ChannelMessageSend(m.ChannelID, "You don't have enough tokens to generate a follow up image.")
			return
		}

		// Generating the response
		placeholder := "Prompt: " + followUpImagePrompt

		// Immediately responding in the 3 second window before the interaciton times out
		msg, err := s.ChannelMessageSendReply(m.ChannelID, placeholder, m.Message.Reference())
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong. Follow up image was not generated.")
			return
		}

		// Going out to make the OpenAI call to get the proper response
		returnFile, err := external.GetDalleFollowupResponse(followUpImagePrompt, originalImageId)

		if err != nil {
			errString := fmt.Sprintf("%s", err.Error())

			// Not 100% sure this is the approach I want to take with handling errors from the API
			_, err := s.ChannelMessageEdit(m.ChannelID, msg.ID, errString)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Something went wrong. Send help.")
			}

			return
		}

		user.ImageTokens -= (&commands.Image{}).CommandCost()
		persistance.UpdateUser(*user)

		logging.LogEvent(eventType.COMMAND_IMAGE, m.Author.ID, followUpImagePrompt, m.GuildID)

		messageEdit := discordgo.NewMessageEdit(msg.ChannelID, msg.ID)
		messageEdit.Files = []*discordgo.File{&returnFile}
		_, err = s.ChannelMessageEditComplex(messageEdit)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong. Send help.")
		}

		return
	}
}

// @Deprecated
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
