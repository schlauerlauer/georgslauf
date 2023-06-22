package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feed struct {
	ID			int64			`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time		`gorm:"default:CURRENT_TIMESTAMP;" json:"CreatedAt"`
	UpdatedAt	time.Time		`gorm:"default:CURRENT_TIMESTAMP;" json:"UpdatedAt"`
	DeletedAt	gorm.DeletedAt	`json:"DeletedAt"`
	TribeID		int64			`json:"TribeID" gorm:"foreignKey:TribeID;"`
	Tribe		Tribe
	StationID	int64			`json:"StationID" gorm:"foreignKey:StationID;"`
	Station		Station
	Official	bool			`json:"official" gorm:"index"`
	Public		bool			`json:"public" gorm:"index"`
	ImageID		uuid.UUID		`gorm:"foreignKey:ImageID;"`
	Image		Image
	// viewable once?
	FeedViewTS	[]FeedViewTS	`gorm:"foreignKey:ID"`
}

type FeedViewTS struct {
	ID			uuid.UUID		`gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	FeedViewed	time.Time
}
