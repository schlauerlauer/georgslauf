package models

import (
	"time"
	"github.com/google/uuid"
)

type Image struct {
	ID			uuid.UUID	`gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt	time.Time
	// TODO thumbhash thumbnail
}
