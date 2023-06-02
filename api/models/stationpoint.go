package models

import (
	"time"
)

// Points given to a station
type StationPoint struct {
	ID			int64		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	GroupID		int64		`json:"GroupID" gorm:"foreignKey:GroupID;index:idx_sp,unique"`
	StationID	int64		`json:"StationID" gorm:"foreignKey:StationID;index:idx_sp,unique"`
	Value		int64		`json:"value"`
}
