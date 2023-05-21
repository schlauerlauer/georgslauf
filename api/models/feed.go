package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feed struct {
	ID			uint			`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time		`json:"CreatedAt"`
	UpdatedAt	time.Time		`json:"UpdatedAt"`
	DeletedAt	gorm.DeletedAt	`json:"DeletedAt"`
	TribeID		uint			`json:"TribeID" gorm:"foreignKey:TribeID;"`
	Tribe		Tribe
	StationID	uint			`json:"StationID" gorm:"foreignKey:StationID;"`
	Station		Station
	Official	bool			`json:"official" gorm:"index"`
	Public		bool			`json:"public" gorm:"index"`
	ImageID		uuid.UUID		`gorm:"foreignKey:ImageID;"`
	Image		Image
	// viewable once?
}
