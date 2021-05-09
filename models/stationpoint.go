package models

import (
	"time"
)

type StationPoint struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	FromID 		uint		`json:"from"`
	From		Group
	ToID		uint		`json:"to"`
	To			Station
	Value		uint		`json:"value"`
}
