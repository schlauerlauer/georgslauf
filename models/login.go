package models

import (
    "time"
)

type Login struct {
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    LastLogin   time.Time   `json:"lastLogin"`
    Username    string      `json:"username" gorm:"unique;not null"`
    Password    string      `json:"password" gorm:"not null"`
    Reset       bool        `json:"reset"`
    Active      bool        `json:"active"`
    Confirmed   bool        `json:"confirmed"`
    Phone       string      `json:"phone"`
    Email       string      `json:"email"`
    Contact     string      `json:"contact"`
    Avatar      string      `json:"avatar"`
    Permissions string      `json:"permissions"`
    UpdatePW    bool        `json:"updatepw" gorm:"-"`
}
