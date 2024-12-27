package database

import (
	"context"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/go-redis/redis/v8"
)

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Env.Redis.Address,
		Password: config.Env.Redis.Password,
		DB:       config.Env.Redis.DB,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
