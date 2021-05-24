package models

import (
	//"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "github.com/sirupsen/logrus"
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
		log.Fatal("Failed to connect to database!")
	}

	db.AutoMigrate(
		&Group{},
		&GroupPoint{},
		&Role{},
		&Login{},
		&Station{},
		&StationPoint{},
		&Tribe{},
		&Grouping{},
		&Content{},
	)

	DB = db

	log.Info("Database migration sucessful.")
}