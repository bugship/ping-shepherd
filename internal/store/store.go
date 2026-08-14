// Package store persists targets and check results.
package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB is a Postgres-backed store.
type DB struct {
	sql *sql.DB
}

// Open connects to Postgres and pings it.
func Open(ctx context.Context, databaseURL string) (*DB, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}
	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(5)
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return &DB{sql: sqlDB}, nil
}

// Ping reports whether the database is reachable.
func (d *DB) Ping(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return fmt.Errorf("no database")
	}
	return d.sql.PingContext(ctx)
}

// Close releases the database connection.
func (d *DB) Close() {
	if d != nil && d.sql != nil {
		_ = d.sql.Close()
	}
}
