package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func FromEnv() Config {
	return Config{
		Port:        envOr("PORT", "8080"),
		DatabaseURL: strings.TrimSpace(os.Getenv("DATABASE_URL")),
	}
}

func (c Config) Addr() string {
	return ":" + c.Port
}

func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
