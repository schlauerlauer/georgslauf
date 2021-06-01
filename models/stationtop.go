package models

type StationTop struct {
	ID			uint		`json:"id"`
	Station		string		`json:"station"`
	Tribe		string		`json:"tribe"`
	Sum			uint		`json:"sum"`
}

func (StationTop) TableName() string {
	return "station_top"
}
