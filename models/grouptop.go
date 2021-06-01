package models

type GroupTop struct {
	ID			uint		`json:"id"`
	Group		string		`json:"group"`
	Grouping	string		`json:"grouping"`
	Tribe		string		`json:"tribe"`
	Sum			uint		`json:"sum"`
}

func (GroupTop) TableName() string {
	return "group_top"
}
