package persistance

import (
	"context"
	"fmt"
	"strings"

	"zombiezen.com/go/sqlite/sqlitex"
)

func SaveConversationMessage(message ConversationMessage) error {
	if message.ThreadId == "" {
		return fmt.Errorf("thread id is required")
	}
	if message.CommandName == "" {
		return fmt.Errorf("command name is required")
	}
	if message.Role == "" {
		return fmt.Errorf("role is required")
	}
	if strings.TrimSpace(message.Content) == "" {
		return fmt.Errorf("content is required")
	}

	db, err := GetConn(context.Background())
	if err != nil {
		return err
	}
	defer PutConn(db)
	return sqlitex.Execute(
		db,
		`INSERT INTO ConversationMessages
		(ThreadId, MessageId, ParentMessageId, ChannelId, GuildId, CommandName, Role, Content)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		&sqlitex.ExecOptions{
			Args: []any{
				message.ThreadId,
				nullIfEmpty(message.MessageId),
				nullIfEmpty(message.ParentMessageId),
				nullIfEmpty(message.ChannelId),
				nullIfEmpty(message.GuildId),
				message.CommandName,
				message.Role,
				message.Content,
			},
		},
	)
}

func GetConversationThreadByMessageID(messageID string, maxMessages int) (*ConversationThread, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, nil
	}

	db, err := GetConn(context.Background())
	if err != nil {
		return nil, err
	}
	defer PutConn(db)

	stmt, err := db.Prepare(`
		SELECT ThreadId, CommandName
		FROM ConversationMessages
		WHERE MessageId = ?
		ORDER BY ID DESC
		LIMIT 1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Finalize()

	stmt.BindText(1, messageID)
	hasRow, err := stmt.Step()
	if err != nil {
		return nil, err
	}
	if !hasRow {
		return nil, nil
	}

	threadID := stmt.ColumnText(0)
	commandName := stmt.ColumnText(1)

	if maxMessages <= 0 {
		maxMessages = 30
	}

	threadMessages, err := getConversationMessagesByThreadID(threadID, maxMessages)
	if err != nil {
		return nil, err
	}

	return &ConversationThread{
		ThreadId:    threadID,
		CommandName: commandName,
		Messages:    threadMessages,
	}, nil
}

func GetConversationMessageByMessageID(messageID string) (*ConversationMessage, error) {
	if strings.TrimSpace(messageID) == "" {
		return nil, nil
	}

	var matchCount int64
	countQuery := `
		SELECT COUNT(*)
		FROM ConversationMessages
		WHERE MessageId = ?
	`
	err := RunQuery(countQuery, &matchCount, messageID)
	if err != nil {
		return nil, err
	}
	if matchCount == 0 {
		return nil, nil
	}

	message := &ConversationMessage{}
	selectQuery := `
		SELECT ID, ThreadId, MessageId, ParentMessageId, ChannelId, GuildId, CommandName, Role, Content
		FROM ConversationMessages
		WHERE MessageId = ?
		ORDER BY ID DESC
		LIMIT 1
	`
	err = RunQuery(selectQuery, message, messageID)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func getConversationMessagesByThreadID(threadID string, maxMessages int) ([]ConversationMessage, error) {
	db, err := GetConn(context.Background())
	if err != nil {
		return nil, err
	}
	defer PutConn(db)

	stmt, err := db.Prepare(`
		SELECT ID, ThreadId, MessageId, ParentMessageId, ChannelId, GuildId, CommandName, Role, Content
		FROM (
			SELECT ID, ThreadId, MessageId, ParentMessageId, ChannelId, GuildId, CommandName, Role, Content
			FROM ConversationMessages
			WHERE ThreadId = ?
			ORDER BY ID DESC
			LIMIT ?
		)
		ORDER BY ID ASC
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Finalize()

	stmt.BindText(1, threadID)
	stmt.BindInt64(2, int64(maxMessages))

	messages := make([]ConversationMessage, 0)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}

		messages = append(messages, ConversationMessage{
			ID:              stmt.ColumnInt64(0),
			ThreadId:        stmt.ColumnText(1),
			MessageId:       stmt.ColumnText(2),
			ParentMessageId: stmt.ColumnText(3),
			ChannelId:       stmt.ColumnText(4),
			GuildId:         stmt.ColumnText(5),
			CommandName:     stmt.ColumnText(6),
			Role:            stmt.ColumnText(7),
			Content:         stmt.ColumnText(8),
		})
	}

	return messages, nil
}

func nullIfEmpty(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
