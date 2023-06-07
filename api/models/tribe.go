package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tribe struct {
	ID			int64			`json:"id" gorm:"primary_key"`
	CreatedAt	time.Time		`json:"CreatedAt"`
	UpdatedAt	time.Time		`json:"UpdatedAt"`
	DeletedAt	gorm.DeletedAt	`json:"DeletedAt"`
	Name		string			`json:"name" gorm:"unique;not null"`
	Short		string			`json:"short" gorm:"unique; not null"`
	DPSG		string			`json:"dpsg"`
	Address		string			`json:"address"`
	URL			string			`json:"url"`
	ImageID		uuid.UUID		`gorm:"foreignKey:ImageID"`
	Image		Image
}
