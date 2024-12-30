package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/G-Villarinho/food-shop-api/cache"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/services/email"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_SignIn(t *testing.T) {
	ctx := context.Background()

	t.Run("should send magic link successfully", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		queueService := &mocks.QueueService{}
		userRepository := &mocks.UserRepository{}

		authService := &authService{
			cacheService:    cacheService,
			queueService:    queueService,
			userRespository: userRepository,
			emailFactory:    *email.NewEmailTaskFactory(),
		}

		email := "user@example.com"
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Email:     email,
			FullName:  "Test User",
		}

		userRepository.On("GetUserByEmail", ctx, email).Return(user, nil)
		cacheService.On("Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute).Return(nil)
		queueService.On("Publish", QueueSendEmail, mock.Anything).Return(nil)

		err := authService.SignIn(ctx, email)

		assert.NoError(t, err)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, email)
		cacheService.AssertCalled(t, "Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute)
		queueService.AssertCalled(t, "Publish", QueueSendEmail, mock.Anything)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		queueService := &mocks.QueueService{}
		userRepository := &mocks.UserRepository{}

		authService := &authService{
			cacheService:    cacheService,
			queueService:    queueService,
			userRespository: userRepository,
			emailFactory:    *email.NewEmailTaskFactory(),
		}

		email := "user@example.com"

		userRepository.On("GetUserByEmail", ctx, email).Return(nil, nil)

		err := authService.SignIn(ctx, email)

		assert.ErrorIs(t, err, models.ErrUserNotFound)
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, email)
		cacheService.AssertNotCalled(t, "Set", ctx, mock.Anything, mock.Anything, mock.Anything)
		queueService.AssertNotCalled(t, "Publish", QueueSendEmail, mock.Anything)
	})

	t.Run("should return error when GetUserByEmail fails", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		queueService := &mocks.QueueService{}
		userRepository := &mocks.UserRepository{}

		authService := &authService{
			cacheService:    cacheService,
			queueService:    queueService,
			userRespository: userRepository,
			emailFactory:    *email.NewEmailTaskFactory(),
		}

		email := "user@example.com"

		userRepository.On("GetUserByEmail", ctx, email).Return(nil, errors.New("database error"))

		err := authService.SignIn(ctx, email)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, email)
		cacheService.AssertNotCalled(t, "Set", ctx, mock.Anything, mock.Anything, mock.Anything)
		queueService.AssertNotCalled(t, "Publish", QueueSendEmail, mock.Anything)
	})

	t.Run("should return error when Set in cache fails", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		queueService := &mocks.QueueService{}
		userRepository := &mocks.UserRepository{}

		authService := &authService{
			cacheService:    cacheService,
			queueService:    queueService,
			userRespository: userRepository,
			emailFactory:    *email.NewEmailTaskFactory(),
		}

		email := "user@example.com"
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Email:     email,
			FullName:  "Test User",
		}

		userRepository.On("GetUserByEmail", ctx, email).Return(user, nil)
		cacheService.On("Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute).Return(errors.New("cache error"))

		err := authService.SignIn(ctx, email)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "set magic link")
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, email)
		cacheService.AssertCalled(t, "Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute)
		queueService.AssertNotCalled(t, "Publish", QueueSendEmail, mock.Anything)
	})

	t.Run("should return error when Publish to queue fails", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		queueService := &mocks.QueueService{}
		userRepository := &mocks.UserRepository{}

		authService := &authService{
			cacheService:    cacheService,
			queueService:    queueService,
			userRespository: userRepository,
			emailFactory:    *email.NewEmailTaskFactory(),
		}

		email := "user@example.com"
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Email:     email,
			FullName:  "Test User",
		}

		userRepository.On("GetUserByEmail", ctx, email).Return(user, nil)
		cacheService.On("Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute).Return(nil)
		queueService.On("Publish", QueueSendEmail, mock.Anything).Return(errors.New("queue error"))

		err := authService.SignIn(ctx, email)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "publish email task")
		userRepository.AssertCalled(t, "GetUserByEmail", ctx, email)
		cacheService.AssertCalled(t, "Set", ctx, mock.Anything, user.ID.String(), 15*time.Minute)
		queueService.AssertCalled(t, "Publish", QueueSendEmail, mock.Anything)
	})
}

func TestAuthService_VeryfyMagicLink(t *testing.T) {
	ctx := context.Background()

	t.Run("should verify magic link successfully", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		userRepository := &mocks.UserRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}
		sessionService := &mocks.SessionService{}

		authService := &authService{
			cacheService:         cacheService,
			userRespository:      userRepository,
			restaurantRepository: restaurantRepository,
			sessionService:       sessionService,
		}

		code := uuid.New()
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Role:      models.Customer,
		}
		sessionToken := "session_token"

		cacheService.On("Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID")).Run(func(args mock.Arguments) {
			*(args.Get(2).(*uuid.UUID)) = userID
		}).Return(nil)
		cacheService.On("Delete", ctx, getMagicLinkKey(code)).Return(nil)
		userRepository.On("GetUserByID", ctx, userID).Return(user, nil)
		sessionService.On("CreateSession", ctx, userID, (*uuid.UUID)(nil), user.Role).Return(&models.Session{
			Token: sessionToken,
		}, nil)

		token, err := authService.VeryfyMagicLink(ctx, code)

		assert.NoError(t, err)
		assert.Equal(t, sessionToken, token)

		cacheService.AssertCalled(t, "Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID"))
		cacheService.AssertCalled(t, "Delete", ctx, getMagicLinkKey(code))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
		sessionService.AssertCalled(t, "CreateSession", ctx, userID, (*uuid.UUID)(nil), user.Role)
	})

	t.Run("should return error when magic link is not found in cache", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		userRepository := &mocks.UserRepository{}
		sessionService := &mocks.SessionService{}
		restaurantRepository := &mocks.RestaurantRepository{}

		authService := &authService{
			cacheService:         cacheService,
			userRespository:      userRepository,
			restaurantRepository: restaurantRepository,
			sessionService:       sessionService,
		}

		code := uuid.New()

		cacheService.On("Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID")).Return(cache.ErrCacheMiss)

		token, err := authService.VeryfyMagicLink(ctx, code)

		assert.Empty(t, token)
		assert.ErrorIs(t, err, models.ErrMagicLinkNotFound)

		cacheService.AssertCalled(t, "Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID"))
		cacheService.AssertNotCalled(t, "Delete", ctx, getMagicLinkKey(code))
		userRepository.AssertNotCalled(t, "GetUserByID", ctx, mock.Anything)
		sessionService.AssertNotCalled(t, "CreateSession", ctx, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error when creating session fails", func(t *testing.T) {
		cacheService := &mocks.CacheService{}
		userRepository := &mocks.UserRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}
		sessionService := &mocks.SessionService{}

		authService := &authService{
			cacheService:         cacheService,
			userRespository:      userRepository,
			restaurantRepository: restaurantRepository,
			sessionService:       sessionService,
		}

		code := uuid.New()
		userID := uuid.New()
		user := &models.User{
			BaseModel: models.BaseModel{ID: userID},
			Role:      models.Customer,
		}

		cacheService.On("Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID")).Run(func(args mock.Arguments) {
			*(args.Get(2).(*uuid.UUID)) = userID
		}).Return(nil)
		cacheService.On("Delete", ctx, getMagicLinkKey(code)).Return(nil)
		userRepository.On("GetUserByID", ctx, userID).Return(user, nil)
		sessionService.On("CreateSession", ctx, userID, (*uuid.UUID)(nil), user.Role).Return(nil, errors.New("session creation error"))

		token, err := authService.VeryfyMagicLink(ctx, code)

		assert.Empty(t, token)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create session")

		cacheService.AssertCalled(t, "Get", ctx, getMagicLinkKey(code), mock.AnythingOfType("*uuid.UUID"))
		cacheService.AssertNotCalled(t, "Delete", ctx, getMagicLinkKey(code))
		userRepository.AssertCalled(t, "GetUserByID", ctx, userID)
		sessionService.AssertCalled(t, "CreateSession", ctx, userID, (*uuid.UUID)(nil), user.Role)
	})
}

func TestAuthService_SignOut(t *testing.T) {
	t.Run("should sign out successfully", func(t *testing.T) {
		sessionID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.SessionIDKey, sessionID)

		sessionService := &mocks.SessionService{}
		sessionService.On("DeleteSession", ctx, sessionID).Return(nil)

		authService := &authService{
			sessionService: sessionService,
		}

		err := authService.SignOut(ctx)

		assert.NoError(t, err)
		sessionService.AssertCalled(t, "DeleteSession", ctx, sessionID)
	})

	t.Run("should return error when session ID is not in context", func(t *testing.T) {
		ctx := context.Background()

		sessionService := &mocks.SessionService{}

		authService := &authService{
			sessionService: sessionService,
		}

		err := authService.SignOut(ctx)

		assert.ErrorIs(t, err, models.ErrSessionNotFound)
		sessionService.AssertNotCalled(t, "DeleteSession", ctx, mock.Anything)
	})

	t.Run("should return error when DeleteSession fails", func(t *testing.T) {
		sessionID := uuid.New()
		ctx := context.WithValue(context.Background(), internal.SessionIDKey, sessionID)

		sessionService := &mocks.SessionService{}
		sessionService.On("DeleteSession", ctx, sessionID).Return(errors.New("delete session error"))

		authService := &authService{
			sessionService: sessionService,
		}

		err := authService.SignOut(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete session")
		sessionService.AssertCalled(t, "DeleteSession", ctx, sessionID)
	})
}
