package config

import (
	"log/slog"
	"os"
)

type Config struct {
    Port           string
    DeadlockAPIURL string
    AssetsURL      string
    RedisAddr      string
}

func Load() *Config {
	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		DeadlockAPIURL: getEnv("DEADLOCK_API_URL", "https://api.deadlock-api.com"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		AssetsURL: getEnv("ASSETS_URL", "https://assets.deadlock-api.com"),
	}

	slog.Info("config loaded",
	"port", cfg.Port,
	"deadlock_api_url", cfg.DeadlockAPIURL,
	"redis_addr", cfg.RedisAddr,
	"assets_url", cfg.AssetsURL,
	)

	return cfg
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    slog.Warn("env var not set, using default", "key", key, "default", fallback)
    return fallback
}
