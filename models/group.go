package models

import (
	"time"
)

type Group struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
  	UpdatedAt	time.Time
	Short 		string		`json:"short"`
	Name		string		`json:"name"`
	Size		uint		`json:"size"`
	RoleID		uint		`json:"role"`
	//Role		GroupRole
	TribeID		uint
	Tribe		Tribe		//FIXME`json:"tribe"`
	Details		string		`json:"details"`
	Contact		string		`json:"contact"`
}
