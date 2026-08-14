// Package config loads process settings from the environment.
package config

import (
	"os"
	"strings"
	"time"
)

// Config is runtime settings for the shepherd process.
type Config struct {
	Port          string
	DatabaseURL   string
	CheckInterval time.Duration
	CheckTimeout  time.Duration
}

// FromEnv reads PORT, DATABASE_URL, CHECK_INTERVAL, and CHECK_TIMEOUT.
func FromEnv() Config {
	return Config{
		Port:          envOr("PORT", "8080"),
		DatabaseURL:   strings.TrimSpace(os.Getenv("DATABASE_URL")),
		CheckInterval: durationOr("CHECK_INTERVAL", 30*time.Second),
		CheckTimeout:  durationOr("CHECK_TIMEOUT", 5*time.Second),
	}
}

// Addr is the TCP listen address, including the leading colon.
func (c Config) Addr() string {
	return ":" + c.Port
}

func durationOr(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return fallback
	}
	return d
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
