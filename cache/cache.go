package cache

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

//go:generate mockery --name=CacheService --output=../mocks --outpkg=mocks
type CacheService interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, target any) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	AddToSet(ctx context.Context, key string, value string, ttl time.Duration) error
	RemoveFromSet(ctx context.Context, key string, value string) error
	GetSetMembers(ctx context.Context, key string, target any) error
}
