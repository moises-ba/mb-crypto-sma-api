package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &redisCache{client: client}
}

func (c *redisCache) Set(ctx context.Context, key string, v any, ttl time.Duration) errors.ApiError {
	val, err := json.Marshal(v)
	if err != nil {
		return errors.NewApiError("fail marshall value "+err.Error(), errors.WithError(err), errors.WithKind(errors.Unexpected))
	}

	err = c.client.Set(ctx, key, string(val), ttl).Err()
	if err != nil && err != redis.Nil {
		return errors.NewApiError("fail to create cache on redis "+err.Error(), errors.WithError(err), errors.WithKind(errors.Unexpected))
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, k string, result any) (bool, errors.ApiError) {
	res, err := c.client.Get(ctx, k).Result()
	if err != nil {
		return false, errors.NewApiError("fail to get cache on redis "+err.Error(), errors.WithError(err), errors.WithKind(errors.Unexpected))
	}

	if res != "" {
		if err := json.Unmarshal([]byte(res), result); err != nil {
			return false, errors.NewApiError("fail to unmarshal from redis "+err.Error(), errors.WithError(err), errors.WithKind(errors.Unexpected))
		}
		return true, nil
	}

	return false, nil
}
