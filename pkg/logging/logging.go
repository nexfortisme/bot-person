package logging

import (
	"fmt"
	"log"
	"main/pkg/persistance"
	"time"

	logging "main/pkg/logging/enums"

	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

func LogEvent(eventType logging.EventType, userId string, message string, serverId string) {

	err := persistance.WithConn(nil, func(conn *sqlite.Conn) error {
		return sqlitex.Execute(
			conn,
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
	})
	if err != nil {
		log.Fatalf("Error logging event: %v", err)
	}
}

func GetLatestEvent(userId string, eventType logging.EventType) (Event, error) {

	queryString := fmt.Sprintf("SELECT * FROM Events WHERE EventUser = '%s' AND EventType = %d ORDER BY EventTime DESC LIMIT 1", userId, eventType)

	event := Event{}
	err := persistance.WithConn(nil, func(conn *sqlite.Conn) error {
		stmt, err := conn.Prepare(queryString)
		if err != nil {
			return err
		}
		defer stmt.Finalize()

		hasRow, err := stmt.Step()
		if err != nil {
			return err
		}

		if !hasRow {
			fmt.Println("No row found")
			return fmt.Errorf("no row found")
		}

		parsedTime, err := time.Parse("2006-01-02 15:04:05", stmt.GetText("EventTime"))
		if err != nil {
			return err
		}

		event.ID = stmt.GetText("ID")
		// event.EventType =
		event.EventUser = stmt.GetText("EventUser")
		event.EventData = stmt.GetText("EventData")
		event.EventServer = stmt.GetText("EventServer")
		event.EventTime = parsedTime

		return nil
	})
	if err != nil {
		return Event{}, err
	}

	return event, nil
}

func LogError(err string) {

	insertErr := persistance.WithConn(nil, func(conn *sqlite.Conn) error {
		return sqlitex.Execute(
			conn,
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
	})
	if insertErr != nil {
		log.Fatalf("Error logging Error: %v", insertErr)
	}
}
