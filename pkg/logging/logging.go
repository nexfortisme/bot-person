package logging

import (
	"fmt"
	"log"
	"main/pkg/persistance"
	"time"

	logging "main/pkg/logging/enums"
	models "main/pkg/logging/models"

	"github.com/surrealdb/surrealdb.go"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	db := persistance.GetDB()

	user, _ := persistance.GetUser(userId)

	event := models.Event{}
	event.EventType = eventType
	event.EventTime = time.Now()
	event.EventUser = userId
	event.EventData = message
	event.EventServer = serverId

	createdEvent, _ := db.Create("events", event)

	// Unmarshal data
	marshaledEvent := make([]models.Event, 1)
	err := surrealdb.Unmarshal(createdEvent, &marshaledEvent)
	if err != nil {
		panic(err)
	}

	relateString := fmt.Sprintf("RELATE users:%s->did->events:%s", user.ID, marshaledEvent[0].ID)

	db.Query(relateString, nil)
}

func LogError(err string) {
	log.Fatalf(err)
}
