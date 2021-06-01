package models

import (
	"time"
)

type Content struct {
	ID			uint		`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	Title		string		`json:"title" gorm:"unique"`
	Body		string		`json:"body" gorm:"unique;size:1023"`
}
