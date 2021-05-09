package models

import (
	"time"
)

type Tribe struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Name		string		`json:"name"`
}
