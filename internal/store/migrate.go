package store

import (
	"context"
	"fmt"
)

const schema = `
CREATE TABLE IF NOT EXISTS targets (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS checks (
	id TEXT PRIMARY KEY,
	target_id TEXT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
	up BOOLEAN NOT NULL,
	status_code INTEGER NOT NULL DEFAULT 0,
	latency_ms INTEGER NOT NULL DEFAULT 0,
	error TEXT NOT NULL DEFAULT '',
	checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS checks_target_checked ON checks (target_id, checked_at DESC);
`

// Migrate creates the targets and checks tables if they do not exist.
func (d *DB) Migrate(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return fmt.Errorf("no database")
	}
	_, err := d.sql.ExecContext(ctx, schema)
	return err
}
