package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
  	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Short 		string		`json:"short" binding:"required" gorm:"unique"`
	Name		string		`json:"name" binding:"required"  gorm:"unique"`
	Size		uint		`json:"size" binding:"required"`
	GroupingID	uint		`json:"GroupingID" binding:"required" gorm:"foreignKey:GroupingID;not null"`
	TribeID		uint		`json:"TribeID" binding:"required" gorm:"foreignKey:TribeID;not null"`
}
