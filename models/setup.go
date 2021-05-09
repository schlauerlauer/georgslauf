package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(sqlite.Open("./georgslauf.db"), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Group{})
	db.AutoMigrate(&GroupPoint{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&Login{})
	db.AutoMigrate(&Station{})
	db.AutoMigrate(&StationPoint{})
	db.AutoMigrate(&Tribe{})

	DB = db
}
