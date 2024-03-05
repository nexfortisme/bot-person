package logging

import (
	logging "main/pkg/logging/enums"
	"time"

	"github.com/google/uuid"
)

type LoggingEvent struct {
	tableName     struct{}          `pg:"tbl_bp_event"`
	ID            uuid.UUID         `pg:"type:uuid,default:gen_random_uuid(),pk"` // Set the type to UUID and use PostgreSQL's gen_random_uuid() function for default value
	DateCreated   time.Time         `pg:"date_created, default:CURRENT_TIMESTAMP"`
	EventType     logging.EventType `pg:"event_type"`
	EventValue    string            `pg:"event_value"`
	CreateUser    string            `pg:"create_user"`
	CreateGuild   string            `pg:"create_guild"`
	CreateGuildId string            `pg:"create_guild_id"`
}
