package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache defines the interface for caching operations.
// Interface-based so we can swap in an in-memory mock for tests.
type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

type redisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis-backed cache
func NewRedisCache(addr, password string) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &redisCache{client: client}
}

func (c *redisCache) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss: %s", key)
	}
	if err != nil {
		return fmt.Errorf("redis get: %w", err)
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling value: %w", err)
	}
	return c.client.Set(ctx, key, b, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// TTL constants — centralised so they're easy to tune
const (
	TTLPlayerStats  = 5 * time.Minute
	TTLMatchHistory = 5 * time.Minute
	TTLHeroList     = 1 * time.Hour
	TTLLeaderboard  = 10 * time.Minute
)

// CacheKey helpers — consistent key format prevents typos
func PlayerKey(steamID string) string    { return fmt.Sprintf("player:%s", steamID) }
func MatchesKey(steamID string) string   { return fmt.Sprintf("matches:%s", steamID) }
func HeroStatsKey(steamID string) string { return fmt.Sprintf("herostats:%s", steamID) }
func HeroListKey() string                { return "heroes:list" }
func LeaderboardKey() string             { return "leaderboard" }
