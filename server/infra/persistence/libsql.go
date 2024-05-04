package persistence

import (
	"database/sql"
	"embed"
	"errors"
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
	}

	slog.Info("creating connector")

	connector, err := libsql.NewEmbeddedReplicaConnector(config.Path, config.Remote)
	if err != nil {
		slog.Error("error creating libsql connector", "err", err)
		return nil, errors.New("test")
	}
	slog.Debug("libsql connector created")
	defer connector.Close()

	sqlDb := sql.OpenDB(connector)
	// defer sqlDb.Close() // TODO

	queries := db.New(sqlDb)

	_ = sqlDb.QueryRow("PRAGMA journal_mode = WAL; PRAGMA foreign_keys = on;")

	goose.SetBaseFS(*migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		panic(err)
	}

	if err := goose.Up(sqlDb, "migrations"); err != nil {
		panic(err)
	}

	return &Repository{
		Queries:  queries,
		Location: location,
	}, nil
}
