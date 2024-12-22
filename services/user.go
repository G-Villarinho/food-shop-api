package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/cache"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, payload models.CreateUserPayload) error
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
		return models.ErrUserNotFound
	}

	if err := u.userRepository.CreateUser(ctx, *payload.ToUser()); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}
