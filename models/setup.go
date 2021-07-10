package models

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	log "github.com/sirupsen/logrus"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/casbin/casbin/v2"
)

var (
	DB *gorm.DB
	EN *casbin.Enforcer
)

func ConnectDatabase() {
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
		&Login{},
		&Station{},
		&StationPoint{},
		&Tribe{},
		&Grouping{},
		&Content{},
		&Run{},
		&ContentType{},
		&Config{},
	)
	db.Exec("DROP VIEW group_top")
	db.Exec("DROP VIEW station_top")
	db.Exec("DROP VIEW public_content")
	db.Exec("DROP VIEW station_tribe")
	db.Exec("DROP VIEW group_tribe")
	db.Exec("DROP VIEW station_public")
	groupingTopView := `
		CREATE VIEW group_top AS
		SELECT
			group_id AS 'id',
			groups.name AS 'group',
			groupings.id AS 'grouping_id',
			groupings.name AS 'grouping',
			tribes.id AS 'tribe_id',
			tribes.name AS 'tribe',
			sum(value) as 'sum',
			round(sum(value)/count(value),2) as 'avg'
		FROM group_points
		INNER JOIN groups on groups.id = group_id
		INNER JOIN groupings on groupings.id = groups.grouping_id
		INNER JOIN tribes on tribes.id = groups.tribe_id
		GROUP BY groups.name
		;
	`
	stationTopView := `
		CREATE VIEW station_top AS
		SELECT
			station_id AS 'id',
			stations.name AS 'station',
			tribes.id AS 'tribe_id',
			tribes.name AS 'tribe',
			sum(value) as 'sum',
			round(sum(value)/count(value),2) as 'avg'
		FROM station_points
		INNER JOIN stations on stations.id = station_id
		INNER JOIN tribes on tribes.id = stations.tribe_id
		GROUP BY station_id
		;
	`
	publicContentView := `
		CREATE VIEW public_content AS
		SELECT
			contents.id,
			content_types.name as 'ct',
			sort,
			value
		FROM contents
		INNER JOIN content_types on contenttype_id = content_types.id
		WHERE content_types.public = '1'
		;
	`
	stationTribeView := `
		CREATE VIEW station_tribe AS
		SELECT
			s.id,
			s.created_at,
			s.updated_at,
			s.short,
			s.name as 'station',
			s.size,
			t.name as 'tribe',
			t.login_id as 'tribe_login'
		FROM stations as s
		INNER JOIN tribes as t ON t.id = s.tribe_id
	`
	groupTribeView := `
		CREATE VIEW group_tribe AS
		SELECT
			g.id,
			g.created_at,
			g.updated_at,
			g.short,
			g.name as 'group',
			g.size,
			t.name as 'tribe',
			t.login_id as 'tribe_login'
		FROM groups as g
		INNER JOIN tribes as t ON t.id = g.tribe_id
	`
	stationPublicView := `
		CREATE VIEW station_public AS
		SELECT
			s.id,
			s.name,
			s.short,
			t.name as 'tribe'
		FROM stations as s
		INNER JOIN tribes as t ON t.id = s.tribe_id
	`
	db.Exec(groupingTopView)
	db.Exec(stationTopView)
	db.Exec(publicContentView)
	db.Exec(stationTribeView)
	db.Exec(groupTribeView)
	db.Exec(stationPublicView)
	DB = db

	log.Info("Database migration sucessful.")

}

func SetEnforcer() {
	a, err := gormadapter.NewAdapterByDBWithCustomTable(DB, &Rule{})
	if err != nil {
		log.Fatal("Error: ", err)
	}
	en, err := casbin.NewEnforcer("keymatch_model.conf", a)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	en.LoadPolicy()
	log.Info("Enforcer connected.")

	en.AddPolicy("admin", "/*", "(GET)|(POST)|(PUT)|(DELETE)|(PATCH)")
	en.SavePolicy()

	EN = en
	log.Info("Policies updated.")
}
