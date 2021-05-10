package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"createdat"`
  	UpdatedAt	time.Time	`json:"updatedat"`
	Short 		string		`json:"short"`
	Name		string		`json:"name"`
	Size		uint		`json:"size"`
	RoleID		uint		`json:"roleid"`
	TribeID		uint		`json:"tribeid`
	Details		string		`json:"details"`
	Contact		string		`json:"contact"`
}
