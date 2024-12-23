package services

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/level-up-api/cache"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, payload models.CreateUserPayload) error
	GetUser(ctx context.Context) (*models.UserResponse, error)
}

type userService struct {
	di             *internal.Di
	userRepository repositories.UserRepository
	cacheService   cache.CacheService
}

func NewUserService(di *internal.Di) (UserService, error) {
	userRepository, err := internal.Invoke[repositories.UserRepository](di)
	if err != nil {
		return nil, err
	}

	cacheService, err := internal.Invoke[cache.CacheService](di)
	if err != nil {
		return nil, err
	}

	return &userService{
		di:             di,
		userRepository: userRepository,
		cacheService:   cacheService,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, payload models.CreateUserPayload) error {
	user, err := u.userRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	if user != nil {
		return models.ErrEmailAlreadyExists
	}

	if err := u.userRepository.CreateUser(ctx, *payload.ToUser()); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (u *userService) GetUser(ctx context.Context) (*models.UserResponse, error) {
	userID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, models.ErrUserNotFoundInContext
	}

	var userResponse *models.UserResponse
	err := u.cacheService.Get(ctx, getUserKey(userID), userResponse)
	if err == nil {
		return userResponse, nil
	}

	if err != cache.ErrCacheMiss {
		return nil, fmt.Errorf("get user from cache: %w", err)
	}

	user, err := u.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {
		return nil, models.ErrUserNotFound
	}

	userResponse = user.ToUserResponse()

	ttl := time.Duration(config.Env.Cache.CacheExp) * time.Minute
	if err := u.cacheService.Set(ctx, getUserKey(userID), userResponse, ttl); err != nil {
		return nil, fmt.Errorf("set user to cache: %w", err)
	}

	return userResponse, nil
}

func getUserKey(userID uuid.UUID) string {
	return fmt.Sprintf("user:%s", userID.String())
}
