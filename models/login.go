package models

import (
	"time"
)

type Login struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
  	UpdatedAt	time.Time
	Name		string		`json:"name"`
	Password	string		`json:"password"`
	RoleID		uint		//`json:"roleid"`
	Role		Role
	// TODO fix json
	Salt		string
}
