package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Station struct {
	ID			uint			`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time
	UpdatedAt	time.Time
	DeletedAt	gorm.DeletedAt
	Name		string			`json:"name" gorm:"unique"`
	Short		string			`json:"short" gorm:"unique"`
	Size		uint			`json:"size"`
	TribeID		uint			`json:"TribeID" gorm:"foreignKey:TribeID,index"`
	Tribe		Tribe
	ImageID		uuid.UUID		`gorm:"foreignKey:ImageID"`
	Image		Image
	Latitude	float64
	Longitude	float64
}
