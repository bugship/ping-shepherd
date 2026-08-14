package config

import "testing"

func TestFromEnvDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	cfg := FromEnv()
	if cfg.Port != "8080" {
		t.Fatalf("port: got %q", cfg.Port)
	}
	if cfg.Addr() != ":8080" {
		t.Fatalf("addr: got %q", cfg.Addr())
	}
}

func TestFromEnvOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_URL", "postgres://coop@localhost/canary")
	cfg := FromEnv()
	if cfg.Port != "9090" {
		t.Fatalf("port: got %q", cfg.Port)
	}
	if cfg.DatabaseURL == "" {
		t.Fatal("expected DATABASE_URL")
	}
}
