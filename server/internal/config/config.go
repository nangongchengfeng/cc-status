package config

import "fmt"

const (
	envListenAddr = "CC_USAGE_SERVER_LISTEN_ADDR"
	envSQLitePath = "CC_USAGE_SERVER_SQLITE_PATH"
	envAuthToken  = "CC_USAGE_SERVER_AUTH_TOKEN"
)

// Config 保存 server 首版运行时所需的最小配置。
type Config struct {
	ListenAddr string
	SQLitePath string
	AuthToken  string
}

// Load 从环境变量读取配置，并对必填项做启动前校验。
func Load(envLookup func(string) string) (Config, error) {
	cfg := Config{
		ListenAddr: ":8080",
		SQLitePath: "./server/data/server.db",
		AuthToken:  envLookup(envAuthToken),
	}

	if value := envLookup(envListenAddr); value != "" {
		cfg.ListenAddr = value
	}
	if value := envLookup(envSQLitePath); value != "" {
		cfg.SQLitePath = value
	}
	if cfg.AuthToken == "" {
		return Config{}, fmt.Errorf("%s is required", envAuthToken)
	}

	return cfg, nil
}
