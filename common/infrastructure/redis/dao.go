package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type dbRedis struct {
	Expiration time.Duration
}

func NewDBRedis(expiration int) dbRedis {
	return dbRedis{Expiration: time.Duration(expiration)}
}

func (r dbRedis) Create(
	ctx context.Context, key string, value interface{},
) *redis.StatusCmd {
	return client.Set(ctx, key, value, 500*time.Second)		// TODO to config
}

func (r dbRedis) Get(
	ctx context.Context, key string,
) *redis.StringCmd {
	return client.Get(ctx, key)
}

func (r dbRedis) Delete(
	ctx context.Context, key string,
) *redis.IntCmd {
	return client.Del(ctx, key)
}

func (r dbRedis) Expire(
	ctx context.Context, key string, expire time.Duration,
) *redis.BoolCmd {
	return client.Expire(ctx, key, 3*time.Second)	// TODO to config
}

func DB() *redis.Client {
	return client
}
