package persistence

import (
	"database/sql"
	"georgslauf/domain/db"
	"georgslauf/interfaces/config"
	"log/slog"
	"os"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

type Repository struct {
	Queries  *db.Queries
	Location *time.Location
}

func NewRepository(config *config.DatabaseConfig) (*Repository, error) {
	location, err := time.LoadLocation(config.Timezone)
	if err != nil {
		slog.Error("error parsing timezone", "err", err)
	}

	sqlDb, err := sql.Open("libsql", config.Path)
	if err != nil {
		slog.Error("error opening database", "err", err)
		os.Exit(1)
	}
	// defer sqlDb.Close() // TODO

	queries := db.New(sqlDb)

	_ = sqlDb.QueryRow("PRAGMA journal_mode = WAL; PRAGMA foreign_keys = on;")

	return &Repository{
		Queries:  queries,
		Location: location,
	}, nil
}
