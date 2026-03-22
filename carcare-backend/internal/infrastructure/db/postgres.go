package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"fmt"
)

// NewPostgres opens a PostgreSQL connection using the given DSN.
func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return db, nil
}
