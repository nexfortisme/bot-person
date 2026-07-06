package logging

import (
	"context"
	"fmt"
	"log"
	"main/pkg/persistance"
	"time"

	logging "main/pkg/logging/enums"
)

func GetLatestEvent(userId string, eventType logging.EventType) (Event, error) {

	queryString := fmt.Sprintf("SELECT * FROM Events WHERE EventUser = '%s' AND EventType = %d ORDER BY EventTime DESC LIMIT 1", userId, eventType)

	db, err := persistance.GetConn(context.Background())
	if err != nil {
		return Event{}, err
	}
	defer persistance.PutConn(db)
	event := Event{}

	stmt, err := db.Prepare(queryString)
	if err != nil {
		return Event{}, err
	}
	defer stmt.Finalize()

	hasRow, err := stmt.Step()
	if err != nil {
		return Event{}, err
	}

	if !hasRow {
		return Event{}, fmt.Errorf("no row found")
	} else {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", stmt.GetText("EventTime"))
		if err != nil {
			return Event{}, err
		}

		event.ID = stmt.GetText("ID")
		// event.EventType =
		event.EventUser = stmt.GetText("EventUser")
		event.EventData = stmt.GetText("EventData")
		event.EventServer = stmt.GetText("EventServer")
		event.EventTime = parsedTime
	}

	return event, nil
}

// LogError writes to stdout only. Event logging to the database was removed
// because the synchronous SQLite writes (on a bind mount under Docker) could
// block command handlers past Discord's 3-second interaction window.
func LogError(err string) {
	log.Printf("ERROR: %s", err)
}
