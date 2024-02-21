package models

import (
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectDatabase(config SqlConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.Hostname, config.Username, config.Password, config.Database, config.Port, config.TZ)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		slog.Error("Failed to connect to database!")
		os.Exit(1)
	}

	db.AutoMigrate(
		&Group{},
		&GroupPoint{},
		&Station{},
		&StationPoint{},
		&Tribe{},
		&Config{},
		&Feed{},
		&Image{},
	)

	slog.Info("Database migration sucessful.")

	// db.Exec(`
	//     CREATE OR REPLACE VIEW group_top AS
	//     SELECT
	//         group_id AS 'id',
	//         g.name AS 'group',
	//         groupings.id AS 'grouping_id',
	//         groupings.name AS 'grouping',
	//         tribes.id AS 'tribe_id',
	//         tribes.name AS 'tribe',
	//         sum(value) as 'sum',
	//         round(sum(value)/count(value),2) as 'avg'
	//     FROM group_points
	//     LEFT JOIN ` + "`groups`" + ` g on g.id = group_id
	//     LEFT JOIN groupings on groupings.id = g.grouping_id
	//     LEFT JOIN tribes on tribes.id = g.tribe_id
	//     GROUP BY g.name
	//     ;
	// `);
	// db.Exec(`
	//     CREATE OR REPLACE VIEW station_top AS
	//     SELECT
	//         station_id AS 'id',
	//         stations.name AS 'station',
	//         tribes.id AS 'tribe_id',
	//         tribes.name AS 'tribe',
	//         sum(value) as 'sum',
	//         round(sum(value)/count(value),2) as 'avg'
	//     FROM station_points
	//     LEFT JOIN stations on stations.id = station_id
	//     LEFT JOIN tribes on tribes.id = stations.tribe_id
	//     GROUP BY station_id
	//     ;
	// `);

	DB = db
}
