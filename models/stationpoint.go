package models

import (
	"time"
)

type StationPoint struct { // Points given to a station
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	GroupID 	uint		`json:"GroupID" gorm:"foreignKey:GroupID` // TODO add unique index
	StationID	uint		`json:"StationID" gorm:"foreignKey:StationID`
	Value		uint		`json:"value"`
}
