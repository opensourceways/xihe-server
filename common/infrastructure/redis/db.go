package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func Init(cfg *Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil

}
