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

	"github.com/bwmarrin/discordgo"
)

var badBotRegex = regexp.MustCompile(`(?i)\bbad bot\b`)
var imageFileExtensionRegex = regexp.MustCompile(`(?i)\.(png|jpe?g|gif|webp|bmp|tiff|heic|heif)$`)
var attachmentDownloadClient = &http.Client{Timeout: 20 * time.Second}

const defaultImageOnlyPrompt = "Please describe the attached image(s)."
const maxImageAttachmentBytes = 20 * 1024 * 1024
const streamInProgressSuffix = "\n\n[Response still generating...]"
const longResponseFileSuffix = "\n\n[Full response attached in too-long.txt]"

func ParseMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignoring messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
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

	// TODO - Handle this better. I don't like this and I feel bad about it
	if badBotRegex.MatchString(m.Message.Content) {
		badBotPrompt := fmt.Sprintf("A Discord user named %s said %q. Reply with a short, funny, sarcastic retort about being called a bad bot. Keep it to one sentence.", m.Author.Username, m.Message.Content)
		badBotRetort := external.GetRetortMachineResponse(badBotPrompt, m.Author.ID)
		if strings.TrimSpace(badBotRetort) == "" {
			badBotRetort = util.GetBadBotResponse()
		}
		_, err := s.ChannelMessageSend(m.ChannelID, badBotRetort)
		util.HandleErrors(err)
		return
	} else if strings.HasPrefix(messageContent, "good bot") {
		goodBotRetort := util.GetGoodBotResponse()
		_, err := s.ChannelMessageSend(m.ChannelID, goodBotRetort)
		util.HandleErrors(err)
		return
	} else if strings.HasPrefix(m.Message.Content, "!addTokens") {
		if !util.UserIsAdmin(m.Author.ID) {
			s.ChannelMessageSend(m.ChannelID, "You do not have permissions to run this command")
			return
		}

		targetUserID, tokenCount, err := parseAddTokensCommand(m.Message.Content)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Usage: `!addTokens @user <positive amount>`")
			return
		}

		if persistance.AddBotPersonTokens(tokenCount, targetUserID) {
			s.ChannelMessageSend(m.ChannelID, "Tokens were successfully added.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong. Tokens were not added.")
		}
		return
	} else if strings.HasPrefix(messageContent, ";;lenny") {
		s.ChannelMessageSend(m.ChannelID, "( ͡° ͜ʖ ͡°)")
		return
	}

	// Reply Handling
	if isReply && (replyType == "bot" || replyType == "bot_gpt") {
		validReplyTarget, validationErr := validateBotConversationReplyTarget(s, m)
		if validationErr != nil {
			fmt.Println("Error validating threaded bot reply target:", validationErr)
			return
		}
		if !validReplyTarget {
			return
		}

		threadID, threadIDErr := resolveThreadIDForLock(m.Message.ReferencedMessage.ID)
		if threadIDErr != nil {
			fmt.Println("Error resolving thread ID for lock:", threadIDErr)
			threadID = m.Message.ReferencedMessage.ID
		}

		if !util.TryStartThreadResponse(threadID) {
			_, _ = s.ChannelMessageSendReply(
				m.ChannelID,
				"Please wait for the current response in this thread to finish before replying again.",
				m.Message.Reference(),
			)
			return
		}
		defer util.FinishThreadResponse(threadID)

		err := handleThreadedBotReply(s, m, replyType)
		if err != nil {
			fmt.Println("Error handling threaded bot reply:", err)
			s.ChannelMessageSendReply(m.ChannelID, "Something went wrong while generating a threaded response.", m.Message.Reference())
		}
		return
	} else if isReply && replyType == "image" {
		// Replies to a generated image go through the same automatic
		// edit-or-chat handling as bot mentions, no slash command needed.
		err := handleMentionedBotMessage(s, m)
		if err != nil {
			fmt.Println("Error handling image reply:", err)
			s.ChannelMessageSendReply(m.ChannelID, "Something went wrong while handling the image reply.", m.Message.Reference())
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

func parseAddTokensCommand(content string) (string, int, error) {
	parts := strings.Fields(content)
	if len(parts) != 3 || parts[0] != "!addTokens" {
		return "", 0, fmt.Errorf("invalid addTokens command")
	}

	mention := parts[1]
	if !strings.HasPrefix(mention, "<@") || !strings.HasSuffix(mention, ">") {
		return "", 0, fmt.Errorf("target must be a user mention")
	}
	userID := strings.TrimSuffix(strings.TrimPrefix(mention, "<@"), ">")
	userID = strings.TrimPrefix(userID, "!")
	if userID == "" {
		return "", 0, fmt.Errorf("empty user id")
	}
	for _, character := range userID {
		if character < '0' || character > '9' {
			return "", 0, fmt.Errorf("invalid user id")
		}
	}

	amount, err := strconv.Atoi(parts[2])
	if err != nil || amount <= 0 {
		return "", 0, fmt.Errorf("amount must be a positive integer")
	}
	return userID, amount, nil
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
	message, err := persistance.GetConversationMessageByMessageID(messageID)
	if err != nil {
		return "", err
	}
	if message == nil {
		return "", nil
	}

	if strings.ToLower(strings.TrimSpace(message.Role)) != "assistant" {
		return "", nil
	}

	switch strings.ToLower(strings.TrimSpace(message.CommandName)) {
	case "bot-gpt", "bot_gpt":
		return "bot_gpt", nil
	case "bot":
		return "bot", nil
	case "image":
		return "image", nil
	default:
		return "", nil
	}
}

func resolveThreadIDForLock(messageID string) (string, error) {
	message, err := persistance.GetConversationMessageByMessageID(messageID)
	if err != nil {
		return "", err
	}

	if message != nil && strings.TrimSpace(message.ThreadId) != "" {
		return message.ThreadId, nil
	}

	return messageID, nil
}

func validateBotConversationReplyTarget(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	if m == nil || m.Message == nil || m.Message.ReferencedMessage == nil {
		return false, nil
	}

	referencedMessage := m.Message.ReferencedMessage
	if referencedMessage.Author == nil || referencedMessage.Author.ID != s.State.User.ID {
		return false, nil
	}

	storedMessage, err := persistance.GetConversationMessageByMessageID(referencedMessage.ID)
	if err != nil {
		return false, err
	}

	if storedMessage == nil {
		// Slash-command and reply metadata can still identify a valid target even before persistence completes.
		replyType := detectReplyType(referencedMessage)
		return replyType == "bot" || replyType == "bot_gpt", nil
	}

	if strings.ToLower(strings.TrimSpace(storedMessage.Role)) != "assistant" {
		return false, nil
	}

	switch normalizeCommandName(storedMessage.CommandName) {
	case "bot", "bot-gpt":
		return true, nil
	default:
		return false, nil
	}
}

func handleMentionedBotMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	userMessageContent := extractMentionPrompt(m.Message.Content, s.State.User.ID, m.Mentions)

	// If the message points at an image (attached or replied to) and reads
	// like an edit request, generate an edited image instead of chatting.
	if userMessageContent != "" && external.ImageEditsSupported() {
		if editImageURLs := collectEditableImageURLs(m); len(editImageURLs) > 0 && external.DetectImageEditIntent(userMessageContent) {
			return handleImageEdit(s, m, userMessageContent, editImageURLs)
		}
	}

	imageDataURLs, imageDescriptions, err := encodeImageAttachmentsToDataURLs(collectMessageAndReferencedAttachments(m))
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
	imageDataURLs, imageDescriptions, err := encodeImageAttachmentsToDataURLs(collectMessageAndReferencedAttachments(m))
	if err != nil {
		return err
	}

	if userMessageContent == "" && len(imageDataURLs) == 0 {
		return nil
	}

	modelMessages = append(modelMessages, buildUserChatMessage(userMessageContent, imageDataURLs))

	placeholderContent := appendSuffixWithinDiscordLimit("Thinking...", streamInProgressSuffix)
	responseMessage, err := s.ChannelMessageSendReply(m.ChannelID, placeholderContent, m.Message.Reference())
	if err != nil {
		return err
	}

	var streamFn func(onDelta func(string)) (string, error)
	var fallbackFn func() string
	switch commandName {
	case "bot-gpt":
		streamFn = func(onDelta func(string)) (string, error) {
			return external.StreamOpenAIGPTResponseWithChatMessages(modelMessages, onDelta)
		}
		fallbackFn = func() string {
			return external.GetOpenAIGPTResponseWithChatMessages(modelMessages)
		}
	default:
		streamFn = func(onDelta func(string)) (string, error) {
			return external.StreamLocalLLMResponseWithChatMessages(modelMessages, m.Author.ID, onDelta)
		}
		fallbackFn = func() string {
			return external.GetLocalLLMResponseWithChatMessages(modelMessages, m.Author.ID)
		}
	}

	assistantResponse, streamErr := util.StreamResponseWithThrottledEdits(
		placeholderContent,
		func(assistant string) string {
			trimmedAssistant := strings.TrimSpace(assistant)
			if trimmedAssistant == "" {
				return placeholderContent
			}
			return appendSuffixWithinDiscordLimit(trimmedAssistant, streamInProgressSuffix)
		},
		streamFn,
		func(content string) error {
			_, editErr := s.ChannelMessageEdit(m.ChannelID, responseMessage.ID, content)
			if editErr != nil {
				fmt.Println("Error editing threaded response while streaming:", editErr)
			}
			return editErr
		},
	)

	if streamErr != nil || strings.TrimSpace(assistantResponse) == "" {
		if streamErr != nil {
			fmt.Println("Error streaming threaded response:", streamErr)
		}
		assistantResponse = fallbackFn()
	}

	if strings.TrimSpace(assistantResponse) == "" {
		assistantResponse = "I'm sorry, I don't understand?"
	}

	trimmedAssistantResponse := strings.TrimSpace(assistantResponse)
	if len(trimmedAssistantResponse) > 2000 {
		filePreview := appendSuffixWithinDiscordLimit(trimmedAssistantResponse, longResponseFileSuffix)
		fileObj := util.HandleTooLongResponseWithFileName(trimmedAssistantResponse, "too-long.txt")

		messageEdit := discordgo.NewMessageEdit(m.ChannelID, responseMessage.ID)
		messageEdit.Content = &filePreview
		messageEdit.Files = []*discordgo.File{fileObj}

		_, err = s.ChannelMessageEditComplex(messageEdit)
		if err != nil {
			fmt.Println("Error attaching threaded too-long file:", err)

			_, _ = s.ChannelMessageEdit(m.ChannelID, responseMessage.ID, filePreview)
			_, sendErr := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Content: "Full response attached in too-long.txt",
				Files:   []*discordgo.File{util.HandleTooLongResponseWithFileName(trimmedAssistantResponse, "too-long.txt")},
				Reference: &discordgo.MessageReference{
					MessageID: responseMessage.ID,
					ChannelID: responseMessage.ChannelID,
					GuildID:   responseMessage.GuildID,
				},
			})
			if sendErr != nil {
				fmt.Println("Error sending threaded too-long file fallback:", sendErr)
			}
		}
	} else {
		_, err = s.ChannelMessageEdit(m.ChannelID, responseMessage.ID, trimmedAssistantResponse)
		if err != nil {
			fmt.Println("Error editing final threaded response:", err)
		}
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

// collectMessageAndReferencedAttachments returns the attachments on the
// user's message followed by those on the message it replies to, so images in
// a replied-to message become model context as well.
func collectMessageAndReferencedAttachments(m *discordgo.MessageCreate) []*discordgo.MessageAttachment {
	attachments := make([]*discordgo.MessageAttachment, 0, len(m.Message.Attachments))
	attachments = append(attachments, m.Message.Attachments...)
	if m.Message.ReferencedMessage != nil {
		attachments = append(attachments, m.Message.ReferencedMessage.Attachments...)
	}
	return attachments
}

func collectEditableImageURLs(m *discordgo.MessageCreate) []string {
	var imageURLs []string
	for _, attachment := range collectMessageAndReferencedAttachments(m) {
		if isImageAttachment(attachment) && strings.TrimSpace(attachment.URL) != "" {
			imageURLs = append(imageURLs, attachment.URL)
		}
	}
	return imageURLs
}

func handleImageEdit(s *discordgo.Session, m *discordgo.MessageCreate, prompt string, imageURLs []string) error {
	user, err := persistance.GetUser(m.Author.ID)
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Something went wrong. The image was not edited.", m.Message.Reference())
		return err
	}

	editCost := (&commands.Image{}).CommandCost()
	if user.ImageTokens < editCost {
		s.ChannelMessageSendReply(m.ChannelID, "You don't have enough tokens to edit an image.", m.Message.Reference())
		return nil
	}

	placeholder := "Prompt: " + prompt
	responseMessage, err := s.ChannelMessageSendReply(m.ChannelID, placeholder, m.Message.Reference())
	if err != nil {
		return err
	}

	returnFile, err := external.GetImageEditResponse(prompt, imageURLs)
	if err != nil {
		_, editErr := s.ChannelMessageEdit(m.ChannelID, responseMessage.ID, err.Error())
		if editErr != nil {
			s.ChannelMessageSend(m.ChannelID, "Something went wrong. Send help.")
		}
		return nil
	}

	user.ImageTokens -= editCost
	persistance.UpdateUser(*user)

	messageEdit := discordgo.NewMessageEdit(m.ChannelID, responseMessage.ID)
	messageEdit.Files = []*discordgo.File{&returnFile}
	_, err = s.ChannelMessageEditComplex(messageEdit)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong. Send help.")
	}

	// Persist the edited image response so replying to it (with or without a
	// mention) routes back through the automatic edit-or-chat handling.
	err = persistance.SaveConversationMessage(persistance.ConversationMessage{
		ThreadId:        responseMessage.ID,
		MessageId:       responseMessage.ID,
		ParentMessageId: m.ID,
		ChannelId:       m.ChannelID,
		GuildId:         m.GuildID,
		CommandName:     "image",
		Role:            "assistant",
		Content:         "[edited image] " + prompt,
	})
	if err != nil {
		fmt.Println("Error saving edited image message:", err)
	}

	return nil
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

func appendSuffixWithinDiscordLimit(content string, suffix string) string {
	if len(suffix) >= 2000 {
		return util.TruncateForDiscord(suffix)
	}

	baseBudget := 2000 - len(suffix)
	baseContent := content
	if len(baseContent) > baseBudget {
		if baseBudget > 3 {
			baseContent = baseContent[:baseBudget-3] + "..."
		} else {
			baseContent = baseContent[:baseBudget]
		}
	}

	return baseContent + suffix
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
