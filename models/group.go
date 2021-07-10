package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
  	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Short 		string		`json:"short" binding:"required" gorm:"unique"`
	Name		string		`json:"name" binding:"required"  gorm:"unique;not null"`
	Size		uint		`json:"size" binding:"required"`
	GroupingID	uint		`json:"GroupingID" binding:"required" gorm:"foreignKey:GroupingID;not null"`
	TribeID		uint		`json:"TribeID" binding:"required" gorm:"foreignKey:TribeID;not null"`
}

type GroupTribe struct {
	ID			uint		`json:"id"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Short		string		`json:"short"`
	Group		string		`json:"group"`
	Size		uint		`json:"size"`
	Tribe		string		`json:"tribe"`
	LoginID		uint		`json:"tribe_login"`
}
func (GroupTribe) TableName() string {
	return "group_tribe"
}