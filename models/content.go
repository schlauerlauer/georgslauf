package models

import (
	"time"
)

type Content struct {
	ID				uint		`json:"id" gorm:"primary_key"`
	CreatedAt		time.Time
	UpdatedAt		time.Time
	Title			string		`json:"title" gorm:"unique"`
	Body			string		`json:"body" gorm:"unique;size:1023"`
	RunID			uint		`json:"RunID" gorm:"foreignKey:RunID;not null"`
	ContenttypeID	uint		`json:"ContenttypeID" gorm:"foreignKey:ContenttypeID;not null"`
}
