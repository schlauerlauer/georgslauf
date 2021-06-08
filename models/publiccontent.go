package models

type PublicContent struct {
	ID			uint		`json:"id"`
	CT			string		`json:"ct"`
	Sort		uint		`json:"sort"`
	Value		string		`json:"value"`
}

func (PublicContent) TableName() string {
	return "public_content"
}
