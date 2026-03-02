package persistance

import (
	"context"

	"zombiezen.com/go/sqlite/sqlitex"
)

func SaveLocalLLMLog(entry LocalLLMLog) error {
	db, err := GetConn(context.Background())
	if err != nil {
		return err
	}
	defer PutConn(db)

	return sqlitex.Execute(
		db,
		`INSERT INTO LocalLLMLogs
		(RequestType, UserId, Model, Endpoint, RequestBody, ResponseBody, StatusCode, ErrorMessage)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		&sqlitex.ExecOptions{
			Args: []any{
				entry.RequestType,
				nullIfEmpty(entry.UserId),
				nullIfEmpty(entry.Model),
				entry.Endpoint,
				nullIfEmpty(entry.RequestBody),
				nullIfEmpty(entry.ResponseBody),
				entry.StatusCode,
				nullIfEmpty(entry.ErrorMessage),
			},
		},
	)
}
