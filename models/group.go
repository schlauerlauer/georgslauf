package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
  	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Short 		string		`json:"short"`
	Name		string		`json:"name"`
	Size		uint		`json:"size"`
	RoleID		uint		`json:"RoleID"`
	TribeID		uint		`json:"TribeID`
	Details		string		`json:"details"`
	Contact		string		`json:"contact"`
}
