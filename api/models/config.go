package models

import (
	"time"
	"gorm.io/gorm"
)

type Config struct {
	ID			uint					`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	gorm.DeletedAt
	Name		string					`json:"name" gorm:"uniqueIndex"`
	Value		map[string]interface{}	`gorm:"serializer:json"`
}

/*
config names:
- notice: message
- contact: slack, tel, whatsapp
*/
