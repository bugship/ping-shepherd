package store

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidName = errors.New("name is required")
	ErrInvalidURL  = errors.New("url must be http or https")
)

type Target struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
}

type NewTarget struct {
	Name string
	URL  string
}

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
