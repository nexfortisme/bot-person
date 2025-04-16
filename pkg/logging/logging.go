package logging

import (
	"log"
	"main/pkg/persistance"

	logging "main/pkg/logging/enums"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	event := Event{}

	err := persistance.RunQuery("INSERT INTO events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)", event.EventType.ToInt(), userId, message, serverId)
	if err != nil {
		log.Fatalf("Error logging event: %v", err)
	}

}

func LogError(err string) {
	persistance.RunQuery("INSERT INTO events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)", logging.ERROR, "SYSTEM", err, "SYSTEM")
}
