package models

import (
	"time"
)

// Points given to a group
type GroupPoint struct {
	ID			int64		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP;"`
	StationID	int64		`json:"StationID" gorm:"foreignKey:StationID;index:idx_gp,unique"`
	GroupID		int64		`json:"GroupID" gorm:"foreignKey:GroupID;index:idx_gp,unique"`
	Value		int64		`json:"value"`
}

type PutPoint struct {
	Value		int64		`form:"value"`
}
