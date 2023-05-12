package models

import (
    "time"
)

type Content struct {
    ID              uint        `json:"id" gorm:"primary_key"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
    Sort            uint        `json:"sort"`
    Value           string      `json:"value"`
    RunID           uint        `json:"RunID" gorm:"foreignKey:RunID;not null"`
    ContenttypeID   uint        `json:"ContenttypeID" gorm:"foreignKey:ContenttypeID;not null"`
}

type PublicContent struct {
    ID      uint    `json:"id"`
    CT      string  `json:"ct"`
    Sort    uint    `json:"sort"`
    Value   string  `json:"value"`
}
func (PublicContent) TableName() string {
    return "public_content"
}
