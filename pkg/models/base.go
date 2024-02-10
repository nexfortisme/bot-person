package model

import "time"

type Base struct {
	ID           string `gorm:"primarykey"`
	DateCreated  time.Time
	DateModified time.Time
	CreatedBy    string
	ModifiedBy   string
}
