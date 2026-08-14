package store

import "time"

type Check struct {
	ID         string    `json:"id"`
	TargetID   string    `json:"target_id"`
	Up         bool      `json:"up"`
	StatusCode int       `json:"status_code"`
	LatencyMS  int64     `json:"latency_ms"`
	Error      string    `json:"error,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
}
