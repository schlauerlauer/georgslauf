package models

import (
	"time"
)

type Station struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Name		string		`json:"name" gorm:"unique"`
	Short		string		`json:"short" gorm:"unique"`
	Size		uint		`json:"size"`
	TribeID		uint		`json:"TribeID" gorm:"foreignKey:TribeID"`
	LoginID		uint		`json:"LoginID" gorm:"foreignKey:LoginID"`
}

type StationTribe struct { // Redacted Station table // TODO Add Login Information aswell -> contact, login email and so on
	ID			uint		`json:"id"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Short		string		`json:"short"`
	Station		string		`json:"station"`
	Size		uint		`json:"size"`
	TribeID		string		`json:"TribeID"`
}
func (StationTribe) TableName() string {
	return "stations"
}