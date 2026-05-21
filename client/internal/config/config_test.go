package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigWithEnvOverride(t *testing.T) {
	t.Parallel()

	appDir := filepath.Join(t.TempDir(), ".cc-usage-client")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() returned error: %v", err)
	}

	configPath := filepath.Join(appDir, "config.yaml")
	configContent := []byte("" +
		"server_url: https://example.com\n" +
		"auth_token: from-file\n" +
		"batch_size: 500\n" +
		"timeout_seconds: 30\n")
	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatalf("WriteFile() returned error: %v", err)
	}

	env := map[string]string{
		"CC_USAGE_CLIENT_SERVER_URL":      "https://override.example.com",
		"CC_USAGE_CLIENT_AUTH_TOKEN":      "from-env",
		"CC_USAGE_CLIENT_BATCH_SIZE":      "250",
		"CC_USAGE_CLIENT_TIMEOUT_SECONDS": "45",
	}

	cfg, err := Load(appDir, func(key string) string {
		return env[key]
	})
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.ServerURL != "https://override.example.com" {
		t.Fatalf("unexpected server url: %q", cfg.ServerURL)
	}

	if cfg.AuthToken != "from-env" {
		t.Fatalf("unexpected auth token: %q", cfg.AuthToken)
	}

	if cfg.BatchSize != 250 {
		t.Fatalf("unexpected batch size: %d", cfg.BatchSize)
	}

	if cfg.TimeoutSeconds != 45 {
		t.Fatalf("unexpected timeout seconds: %d", cfg.TimeoutSeconds)
	}
}
