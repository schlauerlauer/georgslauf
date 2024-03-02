package persistence

import (
	"database/sql"
	"georgslauf/domain/db"
	"georgslauf/interfaces/config"
	"log/slog"
	"os"

	_ "github.com/libsql/go-libsql"
)

type Repository struct {
	Queries *db.Queries
}

func NewRepository(config *config.DatabaseConfig) (*Repository, error) {
	sqlDb, err := sql.Open("libsql", config.Path)
	if err != nil {
		slog.Error("error opening database", "err", err)
		os.Exit(1)
	}
	defer sqlDb.Close()

	queries := db.New(sqlDb)

	_, err = sqlDb.Exec("PRAGMA journal_mode = WAL; PRAGMA foreign_keys = on;")
	if err != nil {
		slog.Error("error setting pragma", "err", err)
		os.Exit(1)
	}

	return &Repository{
		Queries: queries,
	}, nil
}
