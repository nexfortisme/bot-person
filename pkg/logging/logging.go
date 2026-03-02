package logging

import (
	"context"
	"fmt"
	"log"
	"main/pkg/persistance"
	"time"

	logging "main/pkg/logging/enums"

	"zombiezen.com/go/sqlite/sqlitex"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	db, err := persistance.GetConn(context.Background())
	if err != nil {
		log.Printf("Error getting db connection for logging event: %v", err)
		return
	}
	defer persistance.PutConn(db)

	err = sqlitex.Execute(
		db,
		"INSERT INTO Events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)",
		&sqlitex.ExecOptions{
			Args: []any{
				eventType,
				userId,
				message, // no escaping needed
				serverId,
			},
		},
	)
	if err != nil {
		log.Printf("Error logging event: %v", err)
	}
}

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

	db, connErr := persistance.GetConn(context.Background())
	if connErr != nil {
		log.Printf("Error getting db connection for logging error: %v", connErr)
		return
	}
	defer persistance.PutConn(db)

	insertErr := sqlitex.Execute(
		db,
		"INSERT INTO Events (EventType, EventUser, EventData, EventServer) VALUES (?, ?, ?, ?)",
		&sqlitex.ExecOptions{
			Args: []any{
				logging.ERROR,
				"SYSTEM",
				err, // no escaping needed
				"SYSTEM",
			},
		},
	)
	if insertErr != nil {
		log.Printf("Error logging Error: %v", insertErr)
	}
}
