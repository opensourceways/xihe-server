package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Create(context.Context, string, interface{}) *redis.StatusCmd
	Get(context.Context, string) *redis.StringCmd
	Delete(context.Context, string) *redis.IntCmd
	Expire(context.Context, string, time.Duration) *redis.BoolCmd
}

func WithContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}
