package messages

import (
	"encoding/base64"
	"fmt"
	"io"
	"main/pkg/external"
	"main/pkg/persistance"
	"main/pkg/util"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"main/pkg/commands"
	"main/pkg/logging"
	eventType "main/pkg/logging/enums"

	"github.com/bwmarrin/discordgo"
)

var badBotRegex = regexp.MustCompile(`(?i)\bbad bot\b`)
var imageFileExtensionRegex = regexp.MustCompile(`(?i)\.(png|jpe?g|gif|webp|bmp|tiff|heic|heif)$`)
var attachmentDownloadClient = &http.Client{Timeout: 20 * time.Second}

const defaultImageOnlyPrompt = "Please describe the attached image(s)."
const maxImageAttachmentBytes = 20 * 1024 * 1024

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
	escapedMessageContent := util.EscapeQuotes(m.Message.Content)
	messageContent := strings.ToLower(escapedMessageContent)

	// persistance.APictureIsWorthAThousand(m.Message.Content, m)

	// TODO - Handle this better. I don't like this and I feel bad about it
	if badBotRegex.MatchString(m.Message.Content) {
		logging.LogEvent(eventType.USER_BAD_BOT, m.Author.ID, "Bad bot command used", m.GuildID)
		badBotPrompt := fmt.Sprintf("A Discord user named %s said %q. Reply with a short, funny, sarcastic retort about being called a bad bot. Keep it to one sentence.", m.Author.Username, m.Message.Content)
		badBotRetort := external.GetRetortMachineResponse(badBotPrompt, m.Author.ID)
		if strings.TrimSpace(badBotRetort) == "" {
			badBotRetort = util.GetBadBotResponse()
		}
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
		return
	} else if strings.HasPrefix(messageContent, "good bot") {
		logging.LogEvent(eventType.USER_GOOD_BOT, m.Author.ID, "Good bot command used", m.GuildID)
		goodBotRetort := util.GetGoodBotResponse()
		_, err := s.ChannelMessageSend(m.ChannelID, goodBotRetort)
		util.HandleErrors(err)
		return
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
		return
	} else if strings.HasPrefix(messageContent, ";;lenny") {
		logging.LogEvent(eventType.LENNY, m.Author.ID, "Lenny command used", m.GuildID)
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
		return
	}

	// Reply Handling
	if isReply && (replyType == "bot" || replyType == "bot_gpt") {
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

	if !mentionsBot(m.Mentions, s.State.User.ID) {
		return
	}

	err := handleMentionedBotMessage(s, m)
	if err != nil {
		fmt.Println("Error handling @BotPerson message:", err)
		s.ChannelMessageSendReply(m.ChannelID, "Something went wrong while generating a response.", m.Message.Reference())
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

func handleMentionedBotMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userMessageContent := extractMentionPrompt(m.Message.Content, s.State.User.ID, m.Mentions)
	imageDataURLs, imageDescriptions, err := encodeImageAttachmentsToDataURLs(m.Message.Attachments)
	if err != nil {
		return err
	}

	if userMessageContent == "" && len(imageDataURLs) == 0 {
		return nil
	}

	modelMessages := []external.OpenAIChatMessage{
		buildUserChatMessage(userMessageContent, imageDataURLs),
	}

	assistantResponse := external.GetLocalLLMResponseWithChatMessages(modelMessages, m.Author.ID)
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

	parentMessageID := ""
	if m.Message.ReferencedMessage != nil {
		parentMessageID = m.Message.ReferencedMessage.ID
	}

	persistedUserContent := buildPersistenceContent(userMessageContent, imageDescriptions)
	logging.LogEvent(eventType.COMMAND_BOT, m.Author.ID, persistedUserContent, m.GuildID)
	logging.LogEvent(eventType.EXTERNAL_GPT_RESPONSE, m.Author.ID, assistantResponse, m.GuildID)

	threadID := responseMessage.ID
	if threadID == "" {
		threadID = m.ID
	}

	err = persistance.SaveConversationMessage(persistance.ConversationMessage{
		ThreadId:        threadID,
		MessageId:       m.ID,
		ParentMessageId: parentMessageID,
		ChannelId:       m.ChannelID,
		GuildId:         m.GuildID,
		CommandName:     "bot",
		Role:            "user",
		Content:         persistedUserContent,
	})
	if err != nil {
		fmt.Println("Error saving mention user message:", err)
	}

	err = persistance.SaveConversationMessage(persistance.ConversationMessage{
		ThreadId:        threadID,
		MessageId:       responseMessage.ID,
		ParentMessageId: m.ID,
		ChannelId:       m.ChannelID,
		GuildId:         m.GuildID,
		CommandName:     "bot",
		Role:            "assistant",
		Content:         assistantResponse,
	})
	if err != nil {
		fmt.Println("Error saving mention assistant message:", err)
	}

	return nil
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

	modelMessages := make([]external.OpenAIChatMessage, 0, len(thread.Messages)+1)
	for _, historicalMessage := range thread.Messages {
		role := strings.ToLower(historicalMessage.Role)
		if role != "user" && role != "assistant" {
			continue
		}
		if strings.TrimSpace(historicalMessage.Content) == "" {
			continue
		}

		modelMessages = append(modelMessages, external.OpenAIChatMessage{
			Role:    role,
			Content: historicalMessage.Content,
		})
	}

	userMessageContent := strings.TrimSpace(m.Message.Content)
	imageDataURLs, imageDescriptions, err := encodeImageAttachmentsToDataURLs(m.Message.Attachments)
	if err != nil {
		return err
	}

	if userMessageContent == "" && len(imageDataURLs) == 0 {
		return nil
	}

	modelMessages = append(modelMessages, buildUserChatMessage(userMessageContent, imageDataURLs))

	var assistantResponse string
	switch commandName {
	case "bot-gpt":
		assistantResponse = external.GetOpenAIGPTResponseWithChatMessages(modelMessages)
	default:
		assistantResponse = external.GetLocalLLMResponseWithChatMessages(modelMessages, m.Author.ID)
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
		Content:         buildPersistenceContent(userMessageContent, imageDescriptions),
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

func buildUserChatMessage(text string, imageDataURLs []string) external.OpenAIChatMessage {
	trimmedText := strings.TrimSpace(text)
	if len(imageDataURLs) == 0 {
		return external.OpenAIChatMessage{
			Role:    "user",
			Content: trimmedText,
		}
	}

	parts := make([]external.OpenAIChatContentPart, 0, len(imageDataURLs)+1)
	if trimmedText != "" {
		parts = append(parts, external.OpenAIChatContentPart{
			Type: "text",
			Text: trimmedText,
		})
	} else {
		parts = append(parts, external.OpenAIChatContentPart{
			Type: "text",
			Text: defaultImageOnlyPrompt,
		})
	}

	for _, imageDataURL := range imageDataURLs {
		parts = append(parts, external.OpenAIChatContentPart{
			Type: "image_url",
			ImageURL: &external.OpenAIChatImageURL{
				URL: imageDataURL,
			},
		})
	}

	return external.OpenAIChatMessage{
		Role:    "user",
		Content: parts,
	}
}

func buildPersistenceContent(text string, imageDescriptions []string) string {
	trimmedText := strings.TrimSpace(text)
	if len(imageDescriptions) == 0 {
		if trimmedText != "" {
			return trimmedText
		}
		return defaultImageOnlyPrompt
	}

	imageDetails := fmt.Sprintf("[image attachments: %s]", strings.Join(imageDescriptions, ", "))
	if trimmedText == "" {
		return imageDetails
	}
	return trimmedText + "\n" + imageDetails
}

func extractMentionPrompt(content string, botID string, mentions []*discordgo.User) string {
	cleaned := content
	for _, mention := range mentions {
		mentionToken := fmt.Sprintf("<@%s>", mention.ID)
		nicknameMentionToken := fmt.Sprintf("<@!%s>", mention.ID)
		replacement := mention.Username

		if mention.ID == botID {
			replacement = ""
		}

		cleaned = strings.ReplaceAll(cleaned, mentionToken, replacement)
		cleaned = strings.ReplaceAll(cleaned, nicknameMentionToken, replacement)
	}

	return strings.TrimSpace(cleaned)
}

func encodeImageAttachmentsToDataURLs(attachments []*discordgo.MessageAttachment) ([]string, []string, error) {
	imageDataURLs := make([]string, 0, len(attachments))
	imageDescriptions := make([]string, 0, len(attachments))

	for _, attachment := range attachments {
		if !isImageAttachment(attachment) {
			continue
		}

		dataURL, err := attachmentToDataURL(attachment)
		if err != nil {
			return nil, nil, err
		}

		imageDataURLs = append(imageDataURLs, dataURL)
		imageDescriptions = append(imageDescriptions, attachmentDescription(attachment))
	}

	return imageDataURLs, imageDescriptions, nil
}

func isImageAttachment(attachment *discordgo.MessageAttachment) bool {
	if attachment == nil {
		return false
	}

	contentType := normalizeContentType(attachment.ContentType)
	if strings.HasPrefix(contentType, "image/") {
		return true
	}

	fileName := strings.TrimSpace(attachment.Filename)
	if fileName != "" && imageFileExtensionRegex.MatchString(fileName) {
		return true
	}

	if attachment.URL != "" {
		imagePath := attachment.URL
		parsedURL, err := url.Parse(attachment.URL)
		if err == nil {
			imagePath = parsedURL.Path
		}
		if imagePath != "" && imageFileExtensionRegex.MatchString(strings.ToLower(filepath.Ext(imagePath))) {
			return true
		}
	}

	return false
}

func attachmentToDataURL(attachment *discordgo.MessageAttachment) (string, error) {
	if attachment == nil || strings.TrimSpace(attachment.URL) == "" {
		return "", fmt.Errorf("image attachment URL is missing")
	}

	req, err := http.NewRequest(http.MethodGet, attachment.URL, nil)
	if err != nil {
		return "", err
	}

	resp, err := attachmentDownloadClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("failed to download image attachment: %s", resp.Status)
	}

	limitedReader := io.LimitReader(resp.Body, maxImageAttachmentBytes+1)
	data, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", fmt.Errorf("image attachment is empty")
	}
	if len(data) > maxImageAttachmentBytes {
		return "", fmt.Errorf("image attachment exceeds %d bytes", maxImageAttachmentBytes)
	}

	contentType := normalizeContentType(attachment.ContentType)
	if contentType == "" {
		contentType = normalizeContentType(resp.Header.Get("Content-Type"))
	}
	if contentType == "" || !strings.HasPrefix(contentType, "image/") {
		contentType = normalizeContentType(http.DetectContentType(data))
	}
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("attachment %q is not an image", attachmentDescription(attachment))
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", contentType, encoded), nil
}

func normalizeContentType(contentType string) string {
	normalized := strings.TrimSpace(strings.ToLower(contentType))
	if index := strings.Index(normalized, ";"); index >= 0 {
		normalized = strings.TrimSpace(normalized[:index])
	}
	return normalized
}

func attachmentDescription(attachment *discordgo.MessageAttachment) string {
	if attachment == nil {
		return "image"
	}
	if strings.TrimSpace(attachment.Filename) != "" {
		return attachment.Filename
	}
	if strings.TrimSpace(attachment.URL) != "" {
		return attachment.URL
	}
	return "image"
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
