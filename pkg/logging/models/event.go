package logging

import (
	logging "main/pkg/logging/enums"
	"time"
)

type Event struct {
	ID          string            `json:"id,omitempty"`
	EventType   logging.EventType `json:"eventId"`
	EventTime   time.Time         `json:"eventTime"`
	EventUser   string            `json:"eventUser"`
	EventData   string            `json:"eventData"`
	EventServer string            `json:"eventServer"`
}
