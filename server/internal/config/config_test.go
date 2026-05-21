package config

import (
	"strings"
	"testing"
)

func TestLoadRequiresStaticToken(t *testing.T) {
	t.Parallel()

	_, err := Load(func(string) string {
		return ""
	})
	if err == nil {
		t.Fatal("expected missing token to fail")
	}

	if !strings.Contains(err.Error(), "CC_USAGE_SERVER_AUTH_TOKEN") {
		t.Fatalf("expected missing token error, got %v", err)
	}
}

func TestLoadUsesDefaultsAndEnvOverrides(t *testing.T) {
	t.Parallel()

	env := map[string]string{
		"CC_USAGE_SERVER_AUTH_TOKEN":  "secret-token",
		"CC_USAGE_SERVER_LISTEN_ADDR": ":9090",
	}

	cfg, err := Load(func(key string) string {
		return env[key]
	})
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.AuthToken != "secret-token" {
		t.Fatalf("unexpected token: %q", cfg.AuthToken)
	}
	if cfg.ListenAddr != ":9090" {
		t.Fatalf("unexpected listen addr: %q", cfg.ListenAddr)
	}
	if cfg.SQLitePath != "./server/data/server.db" {
		t.Fatalf("unexpected sqlite path: %q", cfg.SQLitePath)
	}
}
