package models

type GroupTop struct {
	ID			uint		`json:"id"`
	Group		string		`json:"group"`
	GroupingID	uint		`json:"grouping_id"`
	Grouping	string		`json:"grouping"`
	TribeID		uint		`json:"tribe_id"`
	Tribe		string		`json:"tribe"`
	Sum			uint		`json:"sum"`
	Avg			float64		`json:"avg"`
}

func (GroupTop) TableName() string {
	return "group_top"
}
