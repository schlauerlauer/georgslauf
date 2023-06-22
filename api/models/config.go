package models

import (
	"time"

	"gorm.io/gorm"
)

type Config struct {
	ID					int64					`json:"id" gorm:"primary_key;"`
	CreatedAt			time.Time				`gorm:"default:CURRENT_TIMESTAMP;"`
	UpdatedAt			time.Time				`gorm:"default:CURRENT_TIMESTAMP;"`
	DeletedAt			gorm.DeletedAt
	Notice				string
	System				SystemConfig			`gorm:"serializer:json;"`
	Contact				ContactConfig			`gorm:"serializer:json;"`
	Groupings			[]string				`gorm:"serializer:json;"`
}

type SystemConfig struct {
	// Stations can edit points
	AllowGroupPoints		bool				`json:"allowGroupPoints"`
	// Public can view stations
	PublicStationsVisible	bool				`json:"publicStations"`
	// The run has started
	RunStarted				bool				`json:"runStarted"`
	// The run has finished
	RunEnded				bool				`json:"runEnded"`
}

type ContactConfig struct {
	Slack				string					`json:"slack"`
	Tel					string					`json:"tel"`
	Whatsapp			string					`json:"whatsapp"`
	StationAmount		int64					`json:"stationAmount"`
}
