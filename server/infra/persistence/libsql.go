package persistence

import (
	"database/sql"
	"embed"
	"georgslauf/domain/db"
	"georgslauf/interfaces/config"
	"log/slog"
	"time"

	"github.com/pressly/goose/v3"
	"github.com/tursodatabase/go-libsql"
)

type Repository struct {
	Queries  *db.Queries
	Location *time.Location
}

func NewRepository(config *config.DatabaseConfig, migrations *embed.FS) (*Repository, error) {
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		slog.Error("error parsing timezone", "err", err)
		return nil, err
	}

	var sqlDb *sql.DB
	if config.EnableReplication {
		slog.Debug("using replication; creating connector")
		connector, err := libsql.NewEmbeddedReplicaConnector(config.Path, config.Remote)
		if err != nil {
			slog.Error("error creating libsqld connector", "err", err)
			return nil, err
		}
		slog.Debug("libsql connector created")
		defer connector.Close()

		sqlDb = sql.OpenDB(connector)
		// defer sqlDb.Close() // TODO

		if err := connector.Sync(); err != nil {
			slog.Error("error syncing to libsqld", "err", err)
			return nil, err
		}
	} else {
		slog.Debug("using local only")
		sqlDb, err = sql.Open("libsql", "file:"+config.Path)
		if err != nil {
			slog.Error("error opening database", "err", err)
			return nil, err
		}
		// defer sqlDb.Close() // TODO
	}

	queries := db.New(sqlDb)

	_ = sqlDb.QueryRow("PRAGMA journal_mode = WAL; PRAGMA foreign_keys = on;")

	goose.SetBaseFS(*migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err := goose.Up(sqlDb, "migrations"); err != nil {
		return nil, err
	}

	return &Repository{
		Queries:  queries,
		Location: location,
	}, nil
}
