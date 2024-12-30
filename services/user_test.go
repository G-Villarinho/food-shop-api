package services

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/food-shop-api/cache"
	"github.com/G-Villarinho/food-shop-api/internal"
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

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()

	// Helper para configurar mocks e serviço
	setup := func() (*userService, *mocks.UserRepository, *mocks.CacheService, *mocks.RestaurantRepository) {
		userRepository := &mocks.UserRepository{}
		cacheService := &mocks.CacheService{}
		restaurantRepository := &mocks.RestaurantRepository{}
		authService := &mocks.AuthService{}

		userService := &userService{
			userRepository:       userRepository,
			cacheService:         cacheService,
			restaurantRepository: restaurantRepository,
			authService:          authService,
		}

		return userService, userRepository, cacheService, restaurantRepository
	}

	t.Run("should return user from cache successfully", func(t *testing.T) {
		userService, _, cacheService, _ := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		expectedResponse := &models.UserResponse{
			ID:    userID.String(),
			Email: "test@example.com",
		}

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Run(func(args mock.Arguments) {
			arg := args.Get(2).(*models.UserResponse)
			*arg = *expectedResponse
		}).Return(nil)

		response, err := userService.GetUser(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
	})

	t.Run("should fetch user from repository when cache is missed", func(t *testing.T) {
		userService, userRepository, cacheService, _ := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Return(cache.ErrCacheMiss)

		user := &models.User{
			BaseModel: models.BaseModel{
				ID: userID,
			},
			Email: "test@example.com",
			Role:  models.Customer,
		}

		userRepository.On("GetUserByID", ctx, userID).Return(user, nil)
		cacheService.On("Set", ctx, getUserKey(userID), mock.Anything, mock.Anything).Return(nil)

		response, err := userService.GetUser(ctx)

		assert.NoError(t, err)
		assert.Equal(t, user.ToUserResponse(), response)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
		cacheService.AssertCalled(t, "Set", ctx, getUserKey(userID), mock.Anything, mock.Anything)
	})

	t.Run("should return error when cache fetch fails with unexpected error", func(t *testing.T) {
		userService, _, cacheService, _ := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Return(errors.New("unexpected cache error"))

		response, err := userService.GetUser(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get user from cache")
		assert.Nil(t, response)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
	})

	t.Run("should return error when user is not found in repository", func(t *testing.T) {
		userService, userRepository, cacheService, _ := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Return(cache.ErrCacheMiss)
		userRepository.On("GetUserByID", ctx, userID).Return(nil, nil)

		response, err := userService.GetUser(ctx)

		assert.ErrorIs(t, err, models.ErrUserNotFound)
		assert.Nil(t, response)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
	})

	t.Run("should return error when user ID is missing in context", func(t *testing.T) {
		userService, _, _, _ := setup()

		response, err := userService.GetUser(ctx)

		assert.ErrorIs(t, err, models.ErrUserNotFoundInContext)
		assert.Nil(t, response)
	})

	t.Run("should return error when repository fetch fails", func(t *testing.T) {
		userService, userRepository, cacheService, _ := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Return(cache.ErrCacheMiss)
		userRepository.On("GetUserByID", ctx, userID).Return(nil, errors.New("repository error"))

		response, err := userService.GetUser(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "get user by id")
		assert.Nil(t, response)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
	})

	t.Run("should return manager user with restaurant name", func(t *testing.T) {
		userService, userRepository, cacheService, restaurantRepository := setup()

		userID := uuid.New()
		ctx := context.WithValue(ctx, internal.UserIDKey, userID)

		cacheService.On("Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse")).Return(cache.ErrCacheMiss)

		user := &models.User{
			BaseModel: models.BaseModel{
				ID: userID,
			},
			Email: "test@example.com",
			Role:  models.Manager,
		}

		restaurant := &models.Restaurant{
			Name: "Test Restaurant",
		}

		userRepository.On("GetUserByID", ctx, userID).Return(user, nil)
		restaurantRepository.On("GetRestaurantByUserID", ctx, userID).Return(restaurant, nil)
		cacheService.On("Set", ctx, getUserKey(userID), mock.Anything, mock.Anything).Return(nil)

		response, err := userService.GetUser(ctx)

		assert.NoError(t, err)
		assert.Equal(t, "Test Restaurant", response.RestaurantName)
		cacheService.AssertCalled(t, "Get", ctx, getUserKey(userID), mock.AnythingOfType("*models.UserResponse"))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
		restaurantRepository.AssertCalled(t, "GetRestaurantByUserID", ctx, userID)
		cacheService.AssertCalled(t, "Set", ctx, getUserKey(userID), mock.Anything, mock.Anything)
	})
}
