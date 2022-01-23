package models

type Rule struct {
    ID      uint    `json:"id" gorm:"primary_key;autoIncrement"`
    Ptype   string  `json:"ptype" gorm:"size:100;uniqueIndex:unique_index"`
    V0      string  `json:"v0" gorm:"size:100;uniqueIndex:unique_index"`
    V1      string  `json:"v1" gorm:"size:100;uniqueIndex:unique_index"`
    V2      string  `json:"v2" gorm:"size:100;uniqueIndex:unique_index"`
    V3      string  `json:"v3" gorm:"size:100;uniqueIndex:unique_index"`
    V4      string  `json:"v4" gorm:"size:100;uniqueIndex:unique_index"`
    V5      string  `json:"v5" gorm:"size:100;uniqueIndex:unique_index"`
}

func (Rule) TableName() string {
    return "casbin_rule"
}
