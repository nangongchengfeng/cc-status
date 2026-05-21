package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServerURL      string `yaml:"server_url"`
	AuthToken      string `yaml:"auth_token"`
	BatchSize      int    `yaml:"batch_size"`
	TimeoutSeconds int    `yaml:"timeout_seconds"`
}

func Load(appDir string, envLookup func(string) string) (Config, error) {
	cfg := Config{
		BatchSize:      500,
		TimeoutSeconds: 30,
	}

	configPath := filepath.Join(appDir, "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}

	overrideString(envLookup, "CC_USAGE_CLIENT_SERVER_URL", &cfg.ServerURL)
	overrideString(envLookup, "CC_USAGE_CLIENT_AUTH_TOKEN", &cfg.AuthToken)
	if err := overrideInt(envLookup, "CC_USAGE_CLIENT_BATCH_SIZE", &cfg.BatchSize); err != nil {
		return Config{}, err
	}
	if err := overrideInt(envLookup, "CC_USAGE_CLIENT_TIMEOUT_SECONDS", &cfg.TimeoutSeconds); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func overrideString(envLookup func(string) string, key string, target *string) {
	if envLookup == nil {
		return
	}
	if value := envLookup(key); value != "" {
		*target = value
	}
}

func overrideInt(envLookup func(string) string, key string, target *int) error {
	if envLookup == nil {
		return nil
	}

	value := envLookup(key)
	if value == "" {
		return nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("parse %s: %w", key, err)
	}

	*target = parsed
	return nil
}
