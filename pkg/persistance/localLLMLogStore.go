package persistance

import "zombiezen.com/go/sqlite/sqlitex"

func SaveLocalLLMLog(entry LocalLLMLog) error {
	db := GetDB()

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
