package messages

import (
	"fmt"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"strconv"
	"strings"

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

	isReply := m.Message.ReferencedMessage != nil
	replyType := "message"
	if isReply {
		replyType = detectReplyType(m.Message.ReferencedMessage)
		if replyType == "message" {
			dbReplyType, err := detectReplyTypeFromThreadStore(m.Message.ReferencedMessage.ID)
			if err != nil {
				fmt.Println("Error detecting reply type from thread store:", err)
			} else if dbReplyType != "" {
				replyType = dbReplyType
			}
		}
	}
	messageContent := strings.ToLower(util.EscapeQuotes(m.Message.Content))

	// persistance.APictureIsWorthAThousand(m.Message.Content, m)

	// TODO - Handle this better. I don't like this and I feel bad about it
	if strings.HasPrefix(messageContent, "bad bot") {
		logging.LogEvent(eventType.USER_BAD_BOT, m.Author.ID, "Bad bot command used", m.GuildID)
		badBotRetort := util.GetBadBotResponse()
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
	} else if strings.HasPrefix(messageContent, "good bot") {
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
	} else if strings.HasPrefix(messageContent, ";;lenny") {
		logging.LogEvent(eventType.LENNY, m.Author.ID, "Lenny command used", m.GuildID)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
	}

	// Reply Handling
	if isReply && replyType == "message" {

		// If it is a message response and doesn't mention the bot, we out
		if !mentionsBot(m.Mentions, s.State.User.ID) {
			return
		}
		return
	} else if isReply && (replyType == "bot" || replyType == "bot_gpt") {
		err := handleThreadedBotReply(s, m, replyType)
		if err != nil {
			fmt.Println("Error handling threaded bot reply:", err)
			s.ChannelMessageSendReply(m.ChannelID, "Something went wrong while generating a threaded response.", m.Message.Reference())
		}
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

func detectReplyType(referencedMessage *discordgo.Message) string {
	if referencedMessage == nil || referencedMessage.Interaction == nil {
		return "message"
	}

	switch strings.ToLower(referencedMessage.Interaction.Name) {
	case "bot":
		return "bot"
	case "bot-gpt", "bot_gpt":
		return "bot_gpt"
	case "image":
		return "image"
	default:
		return "message"
	}
}

func detectReplyTypeFromThreadStore(messageID string) (string, error) {
	thread, err := persistance.GetConversationThreadByMessageID(messageID, 1)
	if err != nil {
		return "", err
	}
	if thread == nil {
		return "", nil
	}

	switch normalizeCommandName(thread.CommandName) {
	case "bot-gpt":
		return "bot_gpt", nil
	case "bot":
		return "bot", nil
	default:
		return "", nil
	}
}

func handleThreadedBotReply(s *discordgo.Session, m *discordgo.MessageCreate, replyType string) error {
	if m.Message.ReferencedMessage == nil {
		return fmt.Errorf("referenced message is required for threaded replies")
	}

	thread, err := persistance.GetConversationThreadByMessageID(m.Message.ReferencedMessage.ID, 40)
	if err != nil {
		return err
	}

	commandName := normalizeCommandName(replyType)
	threadID := m.Message.ReferencedMessage.ID

	if thread != nil {
		if thread.ThreadId != "" {
			threadID = thread.ThreadId
		}
		if thread.CommandName != "" {
			commandName = normalizeCommandName(thread.CommandName)
		}
	} else {
		// Backfill older messages that predate conversation persistence.
		seedMessage := persistance.ConversationMessage{
			ThreadId:    threadID,
			MessageId:   m.Message.ReferencedMessage.ID,
			ChannelId:   m.ChannelID,
			GuildId:     m.GuildID,
			CommandName: commandName,
			Role:        "assistant",
			Content:     m.Message.ReferencedMessage.Content,
		}
		seedErr := persistance.SaveConversationMessage(seedMessage)
		if seedErr != nil {
			fmt.Println("Error seeding fallback conversation message:", seedErr)
		}

		thread = &persistance.ConversationThread{
			ThreadId:    threadID,
			CommandName: commandName,
			Messages:    []persistance.ConversationMessage{seedMessage},
		}
	}

	modelMessages := make([]external.OpenAIGPTMessage, 0, len(thread.Messages)+1)
	for _, historicalMessage := range thread.Messages {
		role := strings.ToLower(historicalMessage.Role)
		if role != "user" && role != "assistant" {
			continue
		}
		if strings.TrimSpace(historicalMessage.Content) == "" {
			continue
		}

		modelMessages = append(modelMessages, external.OpenAIGPTMessage{
			Role:    role,
			Content: historicalMessage.Content,
		})
	}

	userMessageContent := strings.TrimSpace(m.Message.Content)
	if userMessageContent == "" {
		return nil
	}

	modelMessages = append(modelMessages, external.OpenAIGPTMessage{
		Role:    "user",
		Content: userMessageContent,
	})

	var assistantResponse string
	switch commandName {
	case "bot-gpt":
		assistantResponse = external.GetOpenAIGPTResponseWithMessages(modelMessages)
	default:
		assistantResponse = external.GetLocalLLMResponseWithMessages(modelMessages, m.Author.ID)
	}

	if strings.TrimSpace(assistantResponse) == "" {
		assistantResponse = "I'm sorry, I don't understand?"
	}

	responseContent := assistantResponse
	if len(responseContent) > 2000 {
		responseContent = responseContent[:1997] + "..."
	}

	responseMessage, err := s.ChannelMessageSendReply(m.ChannelID, responseContent, m.Message.Reference())
	if err != nil {
		return err
	}

	err = persistance.SaveConversationMessage(persistance.ConversationMessage{
		ThreadId:        threadID,
		MessageId:       m.ID,
		ParentMessageId: m.Message.ReferencedMessage.ID,
		ChannelId:       m.ChannelID,
		GuildId:         m.GuildID,
		CommandName:     commandName,
		Role:            "user",
		Content:         userMessageContent,
	})
	if err != nil {
		fmt.Println("Error saving threaded user message:", err)
	}

	err = persistance.SaveConversationMessage(persistance.ConversationMessage{
		ThreadId:        threadID,
		MessageId:       responseMessage.ID,
		ParentMessageId: m.ID,
		ChannelId:       m.ChannelID,
		GuildId:         m.GuildID,
		CommandName:     commandName,
		Role:            "assistant",
		Content:         assistantResponse,
	})
	if err != nil {
		fmt.Println("Error saving threaded assistant message:", err)
	}

	return nil
}

func normalizeCommandName(commandName string) string {
	switch strings.ToLower(commandName) {
	case "bot_gpt", "bot-gpt":
		return "bot-gpt"
	case "bot":
		return "bot"
	default:
		return "bot"
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
