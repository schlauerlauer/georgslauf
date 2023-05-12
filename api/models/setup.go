package models

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    log "github.com/sirupsen/logrus"
    "fmt"
)

var (
    DB *gorm.DB
)

func ConnectDatabase(config SqlConfig) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", config.Hostname, config.Username, config.Password, config.Database, config.Port, config.TZ)
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
    log.Info("Database migration sucessful.")
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
    // db.Exec(`
    //     CREATE OR REPLACE VIEW public_content AS
    //     SELECT
    //         contents.id,
    //         content_types.name as 'ct',
    //         sort,
    //         value
    //     FROM contents
    //     INNER JOIN content_types on contenttype_id = content_types.id
    //     WHERE content_types.public = '1'
    //     ;
    // `);
    // db.Exec(`
    //     CREATE OR REPLACE VIEW station_tribe AS
    //     SELECT
    //         s.id,
    //         s.created_at,
    //         s.updated_at,
    //         s.short,
    //         s.name as 'station',
    //         s.size,
    //         t.name as 'tribe',
    //         t.login_id as 'tribe_login'
    //     FROM stations as s
    //     INNER JOIN tribes as t ON t.id = s.tribe_id
    // `);
    // db.Exec(`
    //     CREATE OR REPLACE VIEW group_tribe AS
    //     SELECT
    //         g.id,
    //         g.created_at,
    //         g.updated_at,
    //         g.short,
    //         g.name as 'group',
    //         g.size,
    //         t.name as 'tribe',
    //         t.login_id as 'tribe_login'
    //     FROM ` + "`groups`" + ` g
    //     INNER JOIN tribes as t ON t.id = g.tribe_id
    // `);
    // db.Exec(`
    //     CREATE OR REPLACE VIEW station_public AS
    //     SELECT
    //         s.id,
    //         s.name,
    //         s.short,
    //         t.name as 'tribe'
    //     FROM stations as s
    //     INNER JOIN tribes as t ON t.id = s.tribe_id
    // `);
    // db.Exec(`
    //     CREATE OR REPLACE VIEW group_public AS
    //     SELECT
    //         g.id,
    //         g.short,
    //         g.name,
    //         groupings.name AS 'grouping',
    //         t.name as 'tribe'
    //     FROM ` + "`groups`" + ` g
    //     INNER JOIN tribes as t ON t.id = g.tribe_id
    //     INNER JOIN groupings on groupings.id = g.grouping_id
    // `);
    DB = db
    log.Warn("Database view creation skipped!")
}
