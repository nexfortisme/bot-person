package logging

import (
	"log"
	"main/pkg/persistance"

	logging "main/pkg/logging/enums"
	models "main/pkg/logging/models"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	event := models.Event{}

	err := persistance.RunQuery("INSERT INTO events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)", nil, event.EventType.ToInt(), userId, message, serverId)
	if err != nil {
		log.Fatalf("Error logging event: %v", err)
	}

}

func LogError(err string) {
	persistance.RunQuery("INSERT INTO events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)", logging.ERROR, "SYSTEM", err, "SYSTEM")
}
