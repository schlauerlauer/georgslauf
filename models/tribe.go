package models

import (
	"time"
)

type Tribe struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"createdat"`
	UpdatedAt	time.Time	`json:"updatedat"`
	Name		string		`json:"name"`
}
