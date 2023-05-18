package models

import (
	"time"
	"github.com/google/uuid"
)

type Tribe struct {
	ID          uint        `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time   `json:"CreatedAt"`
	UpdatedAt   time.Time   `json:"UpdatedAt"`
	DeletedAt   time.Time   `json:"DeletedAt"`
	Name        string      `json:"name" gorm:"unique;not null"`
	Short       string      `json:"short" gorm:"unique; not null"`
	DPSG        string      `json:"dpsg"`
	Address     string      `json:"address"`
	URL         string      `json:"url"`
	ImageID     uuid.UUID        `gorm:"foreignKey:ImageID"`
	Image       Image
}
