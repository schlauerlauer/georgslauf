package models

import (
    "time"
    "database/sql/driver"
)

type Group struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time   `json:"CreatedAt"`
    UpdatedAt   time.Time   `json:"UpdatedAt"`
    Short       string      `json:"short" binding:"required" gorm:"unique"`
    Name        string      `json:"name" binding:"required"  gorm:"unique;not null"`
    Size        uint        `json:"size" binding:"required"`
    Grouping    Grouping    `json:"grouping"`
    TribeID     uint        `json:"TribeID" binding:"required" gorm:"foreignKey:TribeID;not null"`
}

type Grouping string

const (
    Wös Grouping = "Wös"
    Jupfis Grouping = "Jupfis"
    Pfadis Grouping = "Pfadis"
    Rover Grouping = "Rover"
)

func (p *Grouping) Scan(value interface{}) error {
	*p = Grouping(value.([]byte))
	return nil
}

func (p Grouping) Value() (driver.Value, error) {
	return string(p), nil
}
