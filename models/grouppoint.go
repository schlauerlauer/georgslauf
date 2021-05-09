package models

import (
	"time"
)

type GroupPoint struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	FromID 		uint		`json:"from"`
	From		Station
	ToID		uint		`json:"to"`
	To			Group
	Value		uint		`json:"value"`
}
