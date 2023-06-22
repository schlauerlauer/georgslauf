package models

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID			uuid.UUID	`gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt	time.Time	`gorm:"default:CURRENT_TIMESTAMP;"`
	// TODO thumbhash thumbnail
}
