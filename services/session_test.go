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

func TestSessionService_CreateSession(t *testing.T) {
	t.Run("should create session successfully", func(t *testing.T) {
		tokenService := &mocks.TokenService{}
		cacheService := &mocks.CacheService{}
		sessionService := &sessionService{
			tokenService: tokenService,
			cacheService: cacheService,
		}

		userID := uuid.New()
		restaurantID := uuid.New()
		role := models.Manager
		token := "mock-token"

		tokenService.On("CreateToken", userID, mock.AnythingOfType("uuid.UUID")).Return(token, nil)
		cacheService.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		cacheService.On("AddToSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		session, err := sessionService.CreateSession(context.Background(), userID, &restaurantID, role)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserID)
		assert.Equal(t, restaurantID, *session.RestaurantID)
		assert.Equal(t, role, session.Role)
		assert.Equal(t, token, session.Token)

		tokenService.AssertCalled(t, "CreateToken", userID, mock.AnythingOfType("uuid.UUID"))
		cacheService.AssertCalled(t, "Set", mock.Anything, getSessionKey(session.SessionID), mock.Anything, mock.Anything)
		cacheService.AssertCalled(t, "AddToSet", mock.Anything, getUserSessionsKey(userID), session.SessionID.String(), mock.Anything)
	})

	t.Run("should return error if token creation fails", func(t *testing.T) {
		tokenService := &mocks.TokenService{}
		cacheService := &mocks.CacheService{}
		sessionService := &sessionService{
			tokenService: tokenService,
			cacheService: cacheService,
		}

		userID := uuid.New()
		sessionID := uuid.New()
		restaurantID := uuid.New()
		role := models.Manager

		tokenService.On("CreateToken", userID, mock.AnythingOfType("uuid.UUID")).Return("", errors.New("token creation failed"))

		session, err := sessionService.CreateSession(context.Background(), userID, &restaurantID, role)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "token creation failed")

		tokenService.AssertCalled(t, "CreateToken", userID, mock.AnythingOfType("uuid.UUID"))
		cacheService.AssertNotCalled(t, "Set", mock.Anything, getSessionKey(sessionID), mock.Anything, mock.Anything)
		cacheService.AssertNotCalled(t, "AddToSet", mock.Anything, getUserSessionsKey(userID), sessionID.String(), mock.Anything)
	})

	t.Run("should return error if cache set fails", func(t *testing.T) {
		tokenService := &mocks.TokenService{}
		cacheService := &mocks.CacheService{}
		sessionService := &sessionService{
			tokenService: tokenService,
			cacheService: cacheService,
		}

		userID := uuid.New()
		restaurantID := uuid.New()
		role := models.Manager
		token := "mock-token"

		tokenService.On("CreateToken", userID, mock.AnythingOfType("uuid.UUID")).Return(token, nil)
		cacheService.On(
			"Set",
			mock.Anything,
			mock.MatchedBy(func(key string) bool { return key[:8] == "session:" }),
			mock.AnythingOfType("*models.Session"),
			mock.Anything,
		).Return(errors.New("cache set failed"))

		session, err := sessionService.CreateSession(context.Background(), userID, &restaurantID, role)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "cache set failed")

		tokenService.AssertCalled(t, "CreateToken", userID, mock.AnythingOfType("uuid.UUID"))
		cacheService.AssertCalled(
			t,
			"Set",
			mock.Anything,
			mock.MatchedBy(func(key string) bool { return key[:8] == "session:" }),
			mock.AnythingOfType("*models.Session"),
			mock.Anything,
		)
		cacheService.AssertNotCalled(t, "AddToSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return error if adding session to user set fails", func(t *testing.T) {
		tokenService := &mocks.TokenService{}
		cacheService := &mocks.CacheService{}
		sessionService := &sessionService{
			tokenService: tokenService,
			cacheService: cacheService,
		}

		userID := uuid.New()
		restaurantID := uuid.New()
		role := models.Manager
		token := "mock-token"

		tokenService.On("CreateToken", userID, mock.AnythingOfType("uuid.UUID")).Return(token, nil)
		cacheService.On("Set", mock.Anything, mock.Anything, mock.AnythingOfType("*models.Session"), mock.Anything).Return(nil)
		cacheService.On("AddToSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("add to set failed"))

		session, err := sessionService.CreateSession(context.Background(), userID, &restaurantID, role)

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "add to set failed")

		tokenService.AssertCalled(t, "CreateToken", userID, mock.AnythingOfType("uuid.UUID"))
		cacheService.AssertCalled(t, "Set", mock.Anything, mock.Anything, mock.AnythingOfType("*models.Session"), mock.Anything)
		cacheService.AssertCalled(t, "AddToSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}
