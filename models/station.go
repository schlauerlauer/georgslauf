package models

import (
    "time"
)

type Station struct {
    ID			uint		`json:"id" gorm:"primary_key"`
    CreatedAt	time.Time
    UpdatedAt	time.Time
    Name		string		`json:"name" gorm:"unique"`
    Short		string		`json:"short" gorm:"unique"`
    Size		uint		`json:"size"`
    TribeID		uint		`json:"TribeID" gorm:"foreignKey:TribeID"`
    LoginID		uint		`json:"LoginID" gorm:"foreignKey:LoginID"`
}

type StationPublic struct {
    ID			uint		`json:"id"`
    Name		string		`json:"name"`
    Short		string		`json:"short"`
    Tribe		string		`json:"tribe"`
}
func (StationPublic) TableName() string {
    return "station_public"
}

type StationTribe struct {
    ID			uint		`json:"id"`
    CreatedAt	time.Time	`json:"created_at"`
    UpdatedAt	time.Time	`json:"updated_at"`
    Short		string		`json:"short"`
    Station		string		`json:"station"`
    Size		uint		`json:"size"`
    Tribe		string		`json:"tribe"`
    LoginID		uint		`json:"tribe_login"`
}
func (StationTribe) TableName() string {
    return "station_tribe"
}