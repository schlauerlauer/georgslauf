package models

type StationTop struct {
	ID			uint		`json:"id"`
	Station		string		`json:"station"`
	TribeID		uint		`json:"tribe_id"`
	Tribe		string		`json:"tribe"`
	Sum			uint		`json:"sum"`
	Avg			float64		`json:"avg"`
}

func (StationTop) TableName() string {
	return "station_top"
}
