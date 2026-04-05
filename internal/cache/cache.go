package cache

import (
	"context"
	"time"

	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
)

type Cache interface {
	Set(ctx context.Context, key string, v any, ttl time.Duration) errors.ApiError
	Get(ctx context.Context, k string, result any) (bool, errors.ApiError)
}
