package models

import (
	"time"
)

type Tribe struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Name		string		`json:"name"`
}
