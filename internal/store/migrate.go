package store

import (
	"context"
	"fmt"
)

const targetsSchema = `
CREATE TABLE IF NOT EXISTS targets (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	url TEXT NOT NULL,
	enabled BOOLEAN NOT NULL DEFAULT TRUE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

func (d *DB) Migrate(ctx context.Context) error {
	if d == nil || d.sql == nil {
		return fmt.Errorf("no database")
	}
	_, err := d.sql.ExecContext(ctx, targetsSchema)
	return err
}
