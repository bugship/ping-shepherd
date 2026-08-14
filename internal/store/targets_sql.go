package store

import (
	"context"
	"database/sql"
	"time"
)

// CreateTarget inserts a target.
func (d *DB) CreateTarget(ctx context.Context, in NewTarget) (Target, error) {
	norm, err := NormalizeTarget(in)
	if err != nil {
		return Target{}, err
	}
	t := Target{
		ID:        newID(),
		Name:      norm.Name,
		URL:       norm.URL,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
	}
	_, err = d.sql.ExecContext(ctx,
		`INSERT INTO targets (id, name, url, enabled, created_at) VALUES ($1, $2, $3, $4, $5)`,
		t.ID, t.Name, t.URL, t.Enabled, t.CreatedAt,
	)
	if err != nil {
		return Target{}, err
	}
	return t, nil
}

// ListTargets returns all targets, oldest first.
func (d *DB) ListTargets(ctx context.Context) ([]Target, error) {
	rows, err := d.sql.QueryContext(ctx,
		`SELECT id, name, url, enabled, created_at FROM targets ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Target
	for rows.Next() {
		var t Target
		if err := rows.Scan(&t.ID, &t.Name, &t.URL, &t.Enabled, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	if out == nil {
		out = []Target{}
	}
	return out, rows.Err()
}

// GetTarget returns a target by id.
func (d *DB) GetTarget(ctx context.Context, id string) (Target, error) {
	var t Target
	err := d.sql.QueryRowContext(ctx,
		`SELECT id, name, url, enabled, created_at FROM targets WHERE id = $1`,
		id,
	).Scan(&t.ID, &t.Name, &t.URL, &t.Enabled, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return Target{}, ErrNotFound
	}
	if err != nil {
		return Target{}, err
	}
	return t, nil
}

// DeleteTarget removes a target and its checks.
func (d *DB) DeleteTarget(ctx context.Context, id string) error {
	res, err := d.sql.ExecContext(ctx, `DELETE FROM targets WHERE id = $1`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
