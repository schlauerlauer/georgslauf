package models

import (
    "time"
)

type Group struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time   `json:"CreatedAt"`
    UpdatedAt   time.Time   `json:"UpdatedAt"`
    Short       string      `json:"short" binding:"required" gorm:"unique"`
    Name        string      `json:"name" binding:"required"  gorm:"unique;not null"`
    Size        uint        `json:"size" binding:"required"`
    GroupingID  uint        `json:"grouping" gorm:"index"`
    TribeID     uint        `json:"TribeID" binding:"required" gorm:"foreignKey:TribeID;not null"`
    Tribe       Tribe
}
