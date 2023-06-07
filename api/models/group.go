package models

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID			int64			`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time		`json:"CreatedAt"`
	UpdatedAt	time.Time		`json:"UpdatedAt"`
	DeletedAt	gorm.DeletedAt	`json:"DeletedAt"`
	Short		string			`json:"short" binding:"required" gorm:"unique"`
	Name		string			`json:"name" binding:"required"  gorm:"unique;not null"`
	Size		int64			`json:"size" binding:"required"`
	GroupingID	int64			`json:"grouping" gorm:"index"`
	TribeID		int64			`json:"TribeID" binding:"required" gorm:"foreignKey:TribeID;not null"`
	Tribe		Tribe
}

type GroupWithStationPoints struct {
	Name		string
	Value		int64
	ID			int64
	GroupingID	int64
}
