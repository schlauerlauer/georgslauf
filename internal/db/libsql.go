package db

import (
	"database/sql"

	_ "github.com/tursodatabase/go-libsql"
)

type Repository struct {
	Queries *Queries
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func NewRepository(config *DatabaseConfig) (*Repository, error) {
	sqlDb, err := sql.Open("libsql", "file:"+config.Path)
	if err != nil {
		return nil, err
	}

	// defer sqlDb.Close() // TODO

	queries := New(sqlDb)

	if row := sqlDb.QueryRow("PRAGMA foreign_keys = on;"); row.Err() != nil {
		return nil, err
	}

	return &Repository{
		Queries: queries,
	}, nil
}
