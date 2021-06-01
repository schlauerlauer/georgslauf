package models

import (
	"time"
)

type Tribe struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time	`json:"CreatedAt"`
	UpdatedAt	time.Time	`json:"UpdatedAt"`
	Name		string		`json:"name" gorm:"unique;not null"`
	Short		string		`json:"short" gorm:"unique; not null"`
	DPSG		string		`json:"dpsg"`
	Address		string		`json:"address"`
	LoginID		uint		`json:"LoginID" gorm:"foreignKey:LoginID"`
}
