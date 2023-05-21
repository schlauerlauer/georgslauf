package models

import (
	"time"
)

type GroupPoint struct { // Points given to a group
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	StationID	uint		`json:"StationID" gorm:"foreignKey:StationID;index:idx_gp,unique"`
	GroupID		uint		`json:"GroupID" gorm:"foreignKey:GroupID;index:idx_gp,unique"`
	Value		uint		`json:"value"`
}
