package models

import (
	"time"
)

type Station struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Short		string		`json:"short" gorm:"unique"`
	Name		string		`json:"name" gorm:"unique"`
	Size		uint		`json:"size"`
	TribeID		uint		`json:"TribeID" gorm:"foreignKey:TribeID"`
	// TODO add role?
}
