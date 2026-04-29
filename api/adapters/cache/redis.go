// api/adapters/cache/redis.go
// RedisCache is the production CacheInterface implementation backed by Redis
// via the go-redis/v9 client.
package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"myanimeapi/internal/logger"
)

// RedisCache wraps a go-redis client and implements CacheInterface.
type RedisCache struct {
	client *redis.Client
	log    *logger.Logger
}

// NewRedisCache parses redisURL (e.g. "redis://localhost:6379") and returns a
// RedisCache. Returns an error if the URL is invalid or the initial Ping fails.
func NewRedisCache(ctx context.Context, redisURL string) (*RedisCache, error) {
	log := logger.New()

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	log.WithField("addr", opts.Addr).Info("Redis cache connected")
	return &RedisCache{client: client, log: log}, nil
}

// Get retrieves a value by key. Returns ErrCacheMiss when the key is absent.
func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheMiss
	}
	if err != nil {
		r.log.WithFields(map[string]interface{}{"key": key, "error": err.Error()}).
			Warning("Redis GET error")
		return nil, err
	}
	return val, nil
}

// Set stores value under key with the given TTL.
func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		r.log.WithFields(map[string]interface{}{"key": key, "error": err.Error()}).
			Warning("Redis SET error")
		return err
	}
	return nil
}

// Delete removes one or more exact keys.
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		r.log.WithFields(map[string]interface{}{"keys": keys, "error": err.Error()}).
			Warning("Redis DEL error")
		return err
	}
	return nil
}

// DeleteByPattern scans for keys matching pattern and deletes them.
// Uses SCAN to avoid blocking the Redis server with KEYS.
func (r *RedisCache) DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			r.log.WithFields(map[string]interface{}{"pattern": pattern, "error": err.Error()}).
				Warning("Redis SCAN error during DeleteByPattern")
			return err
		}
		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				r.log.WithFields(map[string]interface{}{"pattern": pattern, "error": err.Error()}).
					Warning("Redis DEL error during DeleteByPattern")
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

// Ping checks connectivity to Redis.
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
