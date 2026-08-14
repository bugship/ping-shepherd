package store

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var (
	// ErrNotFound is returned when a target or check does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalidName is returned when a target name is empty.
	ErrInvalidName = errors.New("name is required")
	// ErrInvalidURL is returned when a target URL is not http or https.
	ErrInvalidURL = errors.New("url must be http or https")
)

// Target is an HTTP endpoint to watch.
type Target struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

// NewTarget is the input for creating a Target.
type NewTarget struct {
	Name string
	URL  string
}

// NormalizeTarget trims fields and checks that URL is http or https.
func NormalizeTarget(in NewTarget) (NewTarget, error) {
	name := strings.TrimSpace(in.Name)
	raw := strings.TrimSpace(in.URL)
	if name == "" {
		return NewTarget{}, ErrInvalidName
	}
	u, err := url.Parse(raw)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return NewTarget{}, ErrInvalidURL
	}
	return NewTarget{Name: name, URL: u.String()}, nil
}
