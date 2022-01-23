package models

import (
    "time"
)

type Grouping struct {
    ID			uint		`json:"id" gorm:"primary_key"`
    CreatedAt	time.Time
    UpdatedAt	time.Time
    Name		string		`json:"name" gorm:"unique"`
    Short		string		`json:"short" gorm:"unique"`
}
