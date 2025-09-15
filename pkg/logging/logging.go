package logging

import (
	"fmt"
	"log"
	"main/pkg/persistance"
	"main/pkg/util"
	"time"

	logging "main/pkg/logging/enums"

	"zombiezen.com/go/sqlite/sqlitex"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	queryString := fmt.Sprintf("INSERT INTO Events (EventType, EventUser, EventData, EventServer) VALUES (%d, '%s', '%s', '%s')", eventType, userId, util.EscapeQuotes(message), serverId)

	db := persistance.GetDB()

	err := sqlitex.Execute(db, queryString, nil)
	if err != nil {
		log.Fatalf("Error logging event: %v", err)
	}
}

func GetLatestEvent(userId string, eventType logging.EventType) (Event, error) {

	queryString := fmt.Sprintf("SELECT * FROM Events WHERE EventUser = '%s' AND EventType = %d ORDER BY EventTime DESC LIMIT 1", userId, eventType)

	db := persistance.GetDB()
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
		fmt.Println("No row found")
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

func LogError(err string) {
	persistance.RunQuery("INSERT INTO Events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)", logging.ERROR, "SYSTEM", err, "SYSTEM")
}
