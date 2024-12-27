package services

import (
	"context"
	"fmt"
	"time"

	"github.com/G-Villarinho/food-shop-api/cache"
	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, payload models.CreateUserPayload, role models.Role) (uuid.UUID, error)
	GetUser(ctx context.Context) (*models.UserResponse, error)
}

type userService struct {
	di                   *internal.Di
	authService          AuthService
	cacheService         cache.CacheService
	restaurantRepository repositories.RestaurantRepository
	userRepository       repositories.UserRepository
}

func NewUserService(di *internal.Di) (UserService, error) {
	authService, err := internal.Invoke[AuthService](di)
	if err != nil {
		return nil, err
	}

	cacheService, err := internal.Invoke[cache.CacheService](di)
	if err != nil {
		return nil, err
	}

	restaurantRepository, err := internal.Invoke[repositories.RestaurantRepository](di)
	if err != nil {
		return nil, err
	}

	userRepository, err := internal.Invoke[repositories.UserRepository](di)
	if err != nil {
		return nil, err
	}

	return &userService{
		di:                   di,
		authService:          authService,
		cacheService:         cacheService,
		restaurantRepository: restaurantRepository,
		userRepository:       userRepository,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, payload models.CreateUserPayload, role models.Role) (uuid.UUID, error) {
	user, err := u.userRepository.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get user by email: %w", err)
	}

	if user != nil {
		return uuid.Nil, models.ErrEmailAlreadyExists
	}

	user = payload.ToUser(role)
	if err := u.userRepository.CreateUser(ctx, *user); err != nil {
		return uuid.Nil, fmt.Errorf("create user: %w", err)
	}

	if err := u.authService.SignIn(ctx, payload.Email); err != nil {
		return uuid.Nil, fmt.Errorf("sign in: %w", err)
	}

	return user.ID, nil
}

func (u *userService) GetUser(ctx context.Context) (*models.UserResponse, error) {
	userID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, models.ErrUserNotFoundInContext
	}

	var userResponse models.UserResponse
	err := u.cacheService.Get(ctx, getUserKey(userID), &userResponse)
	if err == nil {
		return &userResponse, nil
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

	userResponse = *user.ToUserResponse()

	if user.Role == models.Manager {
		restaurant, err := u.restaurantRepository.GetRestaurantByUserID(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("get restaurant id by user id: %w", err)
		}

		if restaurant != nil {
			userResponse.RestaurantName = restaurant.Name
		}
	}

	ttl := time.Duration(config.Env.Cache.CacheExp) * time.Minute
	if err := u.cacheService.Set(ctx, getUserKey(userID), userResponse, ttl); err != nil {
		return nil, fmt.Errorf("set user to cache: %w", err)
	}

	return &userResponse, nil
}

func getUserKey(userID uuid.UUID) string {
	return fmt.Sprintf("user:%s", userID.String())
}
