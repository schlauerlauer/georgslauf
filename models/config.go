package models

import (
    "time"
)

type Config struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Name        string      `json:"name" gorm:"unique"`
    ValueB      bool        `json:"valueb" gorm:"not null"`
}
