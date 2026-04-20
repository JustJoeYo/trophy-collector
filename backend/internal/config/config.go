package config

import (
	"fmt"
	"os"
)

// Config holds all environment-based configuration.
// Nothing is hardcoded — all secrets come from env vars.
type Config struct {
	Port           string
	SteamAPIKey    string
	DeadlockAPIURL string
	RedisAddr      string
	RedisPassword  string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		SteamAPIKey:    getEnv("STEAM_API_KEY", ""),
		DeadlockAPIURL: getEnv("DEADLOCK_API_URL", "https://api.deadlock-api.com"),
		RedisAddr:      getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
	}

	// Steam API key is required
	if cfg.SteamAPIKey == "" {
		return nil, fmt.Errorf("STEAM_API_KEY environment variable is required")
	}

	return cfg, nil
}

// getEnv returns the env var value or a fallback default
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
