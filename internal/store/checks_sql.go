package store

import (
	"context"
	"database/sql"
	"time"
)

func (d *DB) RecordCheck(ctx context.Context, targetID string, c Check) (Check, error) {
	if _, err := d.GetTarget(ctx, targetID); err != nil {
		return Check{}, err
	}
	if c.ID == "" {
		c.ID = newID()
	}
	c.TargetID = targetID
	if c.CheckedAt.IsZero() {
		c.CheckedAt = time.Now().UTC()
	}
	_, err := d.sql.ExecContext(ctx,
		`INSERT INTO checks (id, target_id, up, status_code, latency_ms, error, checked_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		c.ID, c.TargetID, c.Up, c.StatusCode, c.LatencyMS, c.Error, c.CheckedAt,
	)
	if err != nil {
		return Check{}, err
	}
	return c, nil
}

func (d *DB) LatestCheck(ctx context.Context, targetID string) (Check, error) {
	if _, err := d.GetTarget(ctx, targetID); err != nil {
		return Check{}, err
	}
	var c Check
	err := d.sql.QueryRowContext(ctx,
		`SELECT id, target_id, up, status_code, latency_ms, error, checked_at
		 FROM checks WHERE target_id = $1
		 ORDER BY checked_at DESC LIMIT 1`,
		targetID,
	).Scan(&c.ID, &c.TargetID, &c.Up, &c.StatusCode, &c.LatencyMS, &c.Error, &c.CheckedAt)
	if err == sql.ErrNoRows {
		return Check{}, ErrNotFound
	}
	if err != nil {
		return Check{}, err
	}
	return c, nil
}

func (d *DB) ListChecks(ctx context.Context, targetID string, limit int) ([]Check, error) {
	if _, err := d.GetTarget(ctx, targetID); err != nil {
		return nil, err
	}
	if limit <= 0 {
		limit = 50
	}
	rows, err := d.sql.QueryContext(ctx,
		`SELECT id, target_id, up, status_code, latency_ms, error, checked_at
		 FROM checks WHERE target_id = $1
		 ORDER BY checked_at DESC LIMIT $2`,
		targetID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Check
	for rows.Next() {
		var c Check
		if err := rows.Scan(&c.ID, &c.TargetID, &c.Up, &c.StatusCode, &c.LatencyMS, &c.Error, &c.CheckedAt); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	if out == nil {
		out = []Check{}
	}
	return out, rows.Err()
}
