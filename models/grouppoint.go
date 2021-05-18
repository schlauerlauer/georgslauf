package models

import (
	"time"
)

type GroupPoint struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	StationID 	uint		`json:"StationID" gorm:"foreignKey:StationID` //TODO add unique index
	GroupID		uint		`json:"GroupID" gorm:"foreignKey:GroupID`
	Value		uint		`json:"value"`
}
