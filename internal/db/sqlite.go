package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewSqlite(config *DatabaseConfig) (*Repository, error) {
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
		return nil, err
	}

	if row := sqlDb.QueryRow("PRAGMA foreign_keys = on;"); row.Err() != nil {
		return nil, err
	}

	if row := sqlDb.QueryRow("PRAGMA wal_checkpoint(TRUNCATE);"); row.Err() != nil {
		return nil, err
	}

	return &Repository{
		Queries: queries,
	}, nil
}
