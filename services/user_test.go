package services

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*userService, *mocks.AuthService, *mocks.UserRepository) {
	authService := &mocks.AuthService{}
	userRepository := &mocks.UserRepository{}
	cacheService := &mocks.CacheService{}
	restaurantRepository := &mocks.RestaurantRepository{}

	userService := &userService{
		authService:          authService,
		userRepository:       userRepository,
		cacheService:         cacheService,
		restaurantRepository: restaurantRepository,
	}

	return userService, authService, userRepository
}

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a user successfully when payload is valid", func(t *testing.T) {
		userService, authService, userRepository := setup()

		payload := models.CreateUserPayload{
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    nil,
		}
		role := models.Customer

		userRepository.On("GetUserByEmail", ctx, payload.Email).Return(nil, nil)
		userRepository.On("CreateUser", ctx, mock.Anything).Return(nil)
		authService.On("SignIn", ctx, payload.Email).Return(nil)

		userID, err := userService.CreateUser(ctx, payload, role)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, userID)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, payload.Email)
		userRepository.AssertCalled(t, "CreateUser", ctx, mock.Anything)
		authService.AssertCalled(t, "SignIn", ctx, payload.Email)
	})

	t.Run("should return ErrEmailAlreadyExists when email is already registered", func(t *testing.T) {
		userService, _, userRepository := setup()

		payload := models.CreateUserPayload{
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    nil,
		}
		role := models.Customer

		user := &models.User{BaseModel: models.BaseModel{
			ID: uuid.New(),
		}, Email: payload.Email}

		userRepository.On("GetUserByEmail", ctx, payload.Email).Return(user, nil)

		userID, err := userService.CreateUser(ctx, payload, role)

		assert.ErrorIs(t, err, models.ErrEmailAlreadyExists)
		assert.Equal(t, uuid.Nil, userID)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, payload.Email)
		userRepository.AssertNotCalled(t, "CreateUser", ctx, mock.Anything)
	})

	t.Run("should return error when GetUserByEmail fails", func(t *testing.T) {
		userService, _, userRepository := setup()

		payload := models.CreateUserPayload{
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    nil,
		}
		role := models.Customer

		userRepository.On("GetUserByEmail", ctx, payload.Email).Return(nil, errors.New("database error"))

		userID, err := userService.CreateUser(ctx, payload, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get user by email")
		assert.Equal(t, uuid.Nil, userID)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, payload.Email)
		userRepository.AssertNotCalled(t, "CreateUser", ctx, mock.Anything)
	})

	t.Run("should return error when CreateUser fails", func(t *testing.T) {
		userService, _, userRepository := setup()

		payload := models.CreateUserPayload{
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    nil,
		}
		role := models.Customer

		userRepository.On("GetUserByEmail", ctx, payload.Email).Return(nil, nil)
		userRepository.On("CreateUser", ctx, mock.Anything).Return(errors.New("repository error"))

		userID, err := userService.CreateUser(ctx, payload, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create user")
		assert.Equal(t, uuid.Nil, userID)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, payload.Email)
		userRepository.AssertCalled(t, "CreateUser", ctx, mock.Anything)
	})

	t.Run("should return error when SignIn fails", func(t *testing.T) {
		userService, authService, userRepository := setup()

		payload := models.CreateUserPayload{
			FullName: "Test User",
			Email:    "test@example.com",
			Phone:    nil,
		}
		role := models.Customer

		userRepository.On("GetUserByEmail", ctx, payload.Email).Return(nil, nil)
		userRepository.On("CreateUser", ctx, mock.Anything).Return(nil)
		authService.On("SignIn", ctx, payload.Email).Return(errors.New("sign in error"))

		userID, err := userService.CreateUser(ctx, payload, role)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sign in")
		assert.Equal(t, uuid.Nil, userID)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, payload.Email)
		userRepository.AssertCalled(t, "CreateUser", ctx, mock.Anything)
		authService.AssertCalled(t, "SignIn", ctx, payload.Email)
	})
}
