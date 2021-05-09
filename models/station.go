package models

import (
	"time"
)

type Station struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Short		string		`json:"short"`
	Name		string		`json:"name"`
	TribeID		uint		`json:"tribe"`
	Tribe		Tribe
}
