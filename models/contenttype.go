package models

import (
    "time"
)

type ContentType struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    Name        string      `json:"name" gorm:"unique;not null"`
    Public      bool        `json:"public"`
}
