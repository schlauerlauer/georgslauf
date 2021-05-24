package models

import (
	"time"
)

type Login struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
  	UpdatedAt	time.Time
	Username	string		`json:"username" gorm:"unique"`
	Password	string		`json:"password"`
	RoleID		uint		`json:"RoleID" gorm:"foreignKey:RoleID"`
	Salt		string		`json:"salt"`
	Reset		bool		`json:"reset"`
	Active		bool		`json:"active"`
	Confirmed	bool		`json:"confirmed"`
	Phone		string		`json:"phone"`
	Email		string		`json:"email"`
	Contact		string		`json:"contact"`
}
