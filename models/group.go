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
	RoleID		uint		`json:"RoleID" binding:"required" gorm:"foreignKey:RoleID"`
	TribeID		uint		`json:"TribeID" binding:"required" gorm:"foreignKey:TribeID"`
	Details		string		`json:"details" binding:"required"`
	Contact		string		`json:"contact" binding:"required"`
}
