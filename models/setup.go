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
	db.Exec("DROP VIEW group_top")
	db.Exec("DROP VIEW station_top")
	groupingTopQuery := `
		CREATE VIEW group_top AS
		SELECT
			group_id AS 'id',
			groups.name AS 'group',
			groupings.name AS 'grouping',
			tribes.name AS 'tribe',
			sum(value) as 'sum'
		FROM group_points
		INNER JOIN groups on groups.id = group_id
		INNER JOIN groupings on groupings.id = groups.grouping_id
		INNER JOIN tribes on tribes.id = groups.tribe_id
		GROUP BY groups.name
		ORDER BY sum DESC
		;
	`
	stationTopQuery := `
		CREATE VIEW station_top AS
		SELECT
			station_id AS 'id',
			stations.name AS 'station',
			tribes.name AS 'tribe',
			sum(value) as 'sum'
		FROM station_points
		INNER JOIN stations on stations.id = station_id
		INNER JOIN tribes on tribes.id = stations.tribe_id
		GROUP BY station_id
		ORDER BY sum DESC
		;
	`
	db.Exec(groupingTopQuery)
	db.Exec(stationTopQuery)

	DB = db

	log.Info("Database migration sucessful.")
}