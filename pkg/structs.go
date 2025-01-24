package pkg

import "time"

type Base struct {
	ID           string
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    string
	ModifiedBy   string
}
