package cache

import (
	"context"
	"time"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

type redisCache struct {
	di     *internal.Di
	client *redis.Client
}

func NewRedisCache(di *internal.Di) (CacheService, error) {
	client, err := internal.Invoke[*redis.Client](di)
	if err != nil {
		return nil, err
	}

	return &redisCache{
		di:     di,
		client: client,
	}, nil
}

func (r *redisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	JSON, err := jsoniter.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, JSON, ttl).Err()
}

func (r *redisCache) Get(ctx context.Context, key string, target any) error {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}

		return err
	}

	return jsoniter.Unmarshal([]byte(result), target)
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *redisCache) AddToSet(ctx context.Context, key string, value string, ttl time.Duration) error {
	pipe := r.client.TxPipeline()
	pipe.SAdd(ctx, key, value)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisCache) GetSetMembers(ctx context.Context, key string, target any) error {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	data, err := jsoniter.Marshal(members)
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(data, target)
}

func (r *redisCache) RemoveFromSet(ctx context.Context, key string, value string) error {
	return r.client.SRem(ctx, key, value).Err()
}
