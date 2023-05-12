package models

import (
    "time"
)

type Run struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time   `json:"CreatedAt"`
    UpdatedAt   time.Time   `json:"UpdatedAt"`
    Year        uint        `json:"year"`
    Note        string      `json:"note"`
    TribeID     uint        `json:"TribeID" gorm:"foreignKey:TribeID;not null"`
}
