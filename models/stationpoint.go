package models

import (
	"time"
)

type StationPoint struct { // Points given to a station
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	GroupID 	uint		`json:"GroupID" gorm:"foreignKey:GroupID;index:idx_sp,unique""`
	StationID	uint		`json:"StationID" gorm:"foreignKey:StationID;index:idx_sp,unique"`
	Value		uint		`json:"value"`
}

type StationTop struct { // View
	ID			uint		`json:"id"`
	Station		string		`json:"station"`
	TribeID		uint		`json:"tribe_id"`
	Tribe		string		`json:"tribe"`
	Sum			uint		`json:"sum"`
	Avg			float64		`json:"avg"`
}
func (StationTop) TableName() string {
	return "station_top"
}