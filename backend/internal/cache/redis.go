package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
    client *redis.Client
}

func NewRedisCache(addr string) Cache {
    client := redis.NewClient(&redis.Options{
        Addr: addr,
    })
    return &redisCache{client: client}
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
    val, err := r.client.Get(ctx, key).Result()
    if errors.Is(err, redis.Nil) {
        return "", ErrCacheMiss
    }
    if err != nil {
        return "", err
    }
    return val, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
    return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
    return r.client.Del(ctx, key).Err()
}
