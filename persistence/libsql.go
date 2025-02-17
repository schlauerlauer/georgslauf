package persistence

import (
	"database/sql"
	"georgslauf/config"
	"georgslauf/db"

	_ "github.com/tursodatabase/go-libsql"
)

type Repository struct {
	Queries *db.Queries
}

func NewRepository(config *config.DatabaseConfig) (*Repository, error) {
	sqlDb, err := sql.Open("libsql", "file:"+config.Path)
	if err != nil {
		return nil, err
	}

	// defer sqlDb.Close() // TODO

	queries := db.New(sqlDb)

	if row := sqlDb.QueryRow("PRAGMA foreign_keys = on;"); row.Err() != nil {
		return nil, err
	}

	return &Repository{
		Queries: queries,
	}, nil
}
