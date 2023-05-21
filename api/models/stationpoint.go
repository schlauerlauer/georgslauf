package models

import (
	"time"
)

type StationPoint struct { // Points given to a station
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	GroupID		uint		`json:"GroupID" gorm:"foreignKey:GroupID;index:idx_sp,unique"`
	StationID	uint		`json:"StationID" gorm:"foreignKey:StationID;index:idx_sp,unique"`
	Value		uint		`json:"value"`
}
