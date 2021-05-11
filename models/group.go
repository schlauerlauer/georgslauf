package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
  	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Short 		string		`json:"short" binding:"required"`
	Name		string		`json:"name" binding:"required"`
	Size		uint		`json:"size" binding:"required"`
	RoleID		uint		`json:"RoleID" binding:"required"`
	TribeID		uint		`json:"TribeID binding:"required"`
	Details		string		`json:"details" binding:"required"`
	Contact		string		`json:"contact" binding:"required"`
}
