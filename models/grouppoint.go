package models

import (
    "time"
)

type GroupPoint struct { // Points given to a group
    ID          uint        `json:"id" gorm:"primary_key"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    StationID   uint        `json:"StationID" gorm:"foreignKey:StationID;index:idx_gp,unique"`
    GroupID     uint        `json:"GroupID" gorm:"foreignKey:GroupID;index:idx_gp,unique"`
    Value       uint        `json:"value"`
}

type GroupTop struct { // View
    ID          uint        `json:"id"`
    Group       string      `json:"group"`
    GroupingID  uint        `json:"grouping_id"`
    Grouping    string      `json:"grouping"`
    TribeID     uint        `json:"tribe_id"`
    Tribe       string      `json:"tribe"`
    Sum         uint        `json:"sum"`
    Avg         float64     `json:"avg"`
}
func (GroupTop) TableName() string {
    return "group_top"
}
