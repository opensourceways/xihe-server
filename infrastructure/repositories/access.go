package repositories

import (
	"context"
	"time"

	coredis "github.com/opensourceways/xihe-server/common/infrastructure/redis"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
)

type Access interface {
	Insert(key, value string) error
	Get(key string) (string, error)
	Expire(key string, expire int64) error
}

func NewAccessRepo(expireDuration int) Access {
	return &accessRepo{cli: coredis.NewDBRedis(expireDuration)}
}

type accessRepo struct {
	cli redis.RedisClient
}

func (impl *accessRepo) Insert(key, value string) error {
	f := func(ctx context.Context) error {
		cmd := impl.cli.Create(ctx, key, value)
		if cmd.Err() != nil {
			return cmd.Err()
		}

		ok, err := cmd.Result()
		if ok != "ok" {
			return err
		}

		return nil
	}

	return redis.WithContext(f)
}

func (impl *accessRepo) Get(key string) (string, error) {
	var value string

	f := func(ctx context.Context) error {
		cmd := impl.cli.Get(ctx, key)
		if cmd.Err() != nil {
			return cmd.Err()
		}

		value = cmd.Val()

		return nil
	}

	if err := redis.WithContext(f); err != nil {
		return "", err
	}

	return value, nil
}

func (impl *accessRepo) Expire(key string, expire int64) error {
	f := func(ctx context.Context) error {
		cmd := impl.cli.Expire(ctx, key, time.Duration(expire))

		return cmd.Err()
	}

	return redis.WithContext(f)
}
