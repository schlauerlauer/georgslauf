package db

import (
	"database/sql"
	"log/slog"

	// _ "github.com/tursodatabase/go-libsql"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	Queries *Queries
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func NewSqlite(config *DatabaseConfig) (*Repository, error) {
	// sqlDb, err := sql.Open("libsql", "file:"+config.Path)
	sqlDb, err := sql.Open("sqlite3", config.Path)
	if err != nil {
		return nil, err
	}

	// defer sqlDb.Close() // TODO

	queries := New(sqlDb)

	if rows, err := sqlDb.Query("PRAGMA journal_mode = WAL;"); err != nil {
		return nil, err
	} else {
		if err := rows.Close(); err != nil {
			return nil, err
		}
	}

	if row := sqlDb.QueryRow("PRAGMA synchronous = NORMAL;"); row.Err() != nil {
		slog.Error("pragma synchronous", "err", row.Err())
		return nil, row.Err()
	}

	if row := sqlDb.QueryRow("PRAGMA foreign_keys = 1;"); row.Err() != nil {
		slog.Error("pragma foreign_keys", "err", row.Err())
		return nil, row.Err()
	}

	if row := sqlDb.QueryRow("PRAGMA wal_checkpoint(TRUNCATE);"); row.Err() != nil {
		slog.Error("pragma wal_checkpoint", "err", row.Err())
		return nil, row.Err()
	}

	return &Repository{
		Queries: queries,
	}, nil
}
