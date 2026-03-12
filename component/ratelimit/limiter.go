package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter interface {
	IsAllowed(ctx context.Context, clientID, path string, limit int64, window time.Duration) bool
}

type RedisLimiter struct {
	rdb    *redis.Client
	prefix string
}

func NewRedisLimiter(rdb *redis.Client, prefix string) *RedisLimiter {
	if prefix == "" {
		prefix = "rl"
	}

	return &RedisLimiter{
		rdb:    rdb,
		prefix: prefix,
	}
}

func (r *RedisLimiter) IsAllowed(ctx context.Context, clientID, path string, limit int64, window time.Duration) bool {
	if clientID == "" || path == "" || limit <= 0 || window <= 0 {
		return false
	}

	key := r.buildKey(clientID, path, window)

	count, err := r.rdb.Incr(ctx, key).Result()
	if err != nil {
		return false
	}

	if count == 1 {
		if err := r.rdb.Expire(ctx, key, window).Err(); err != nil {
			return false
		}
	}

	return count <= limit
}

func (r *RedisLimiter) buildKey(clientID, path string, window time.Duration) string {
	return fmt.Sprintf("%s:%s:%s:%d", r.prefix, clientID, path, int64(window.Seconds()))
}
