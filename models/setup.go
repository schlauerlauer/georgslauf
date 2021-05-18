package models

import (
	//"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// db, err := gorm.Open(sqlite.Open("./georgslauf.db"), &gorm.Config{
	// 	PrepareStmt: true,
	// })

	dsn := "***REMOVED***?charset=utf8mb4&parseTime=true&loc=Local&tls=skip-verify"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(
		&Group{},
		&GroupPoint{},
		&Role{},
		&Login{},
		&Station{},
		&StationPoint{},
		&Tribe{},
	)

	DB = db
}