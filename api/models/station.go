package models

import (
    "time"
)

type Station struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Name        string      `json:"name" gorm:"unique"`
    Short       string      `json:"short" gorm:"unique"`
    Size        uint        `json:"size"`
    TribeID     uint        `json:"TribeID" gorm:"foreignKey:TribeID,index"`
    Tribe       Tribe
    // LoginID     uint        `json:"LoginID" gorm:"foreignKey:LoginID"`
}
