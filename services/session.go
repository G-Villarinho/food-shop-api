package services

import (
	"context"
	"errors"
	"time"

	"github.com/G-Villarinho/food-shop-api/cache"
	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
)

//go:generate mockery --name=SessionService --output=../mocks --outpkg=mocks
type SessionService interface {
	CreateSession(ctx context.Context, userID uuid.UUID, restaurantID *uuid.UUID, role models.Role) (*models.Session, error)
	GetSessionByToken(ctx context.Context, token string) (*models.Session, error)
	GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error)
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	DeleteAllSessions(ctx context.Context, userID uuid.UUID) error
}

type sessionService struct {
	di           *internal.Di
	tokenService TokenService
	cacheService cache.CacheService
}

func NewSessionService(di *internal.Di) (SessionService, error) {
	tokenService, err := internal.Invoke[TokenService](di)
	if err != nil {
		return nil, err
	}

	cacheService, err := internal.Invoke[cache.CacheService](di)
	if err != nil {
		return nil, err
	}

	return &sessionService{
		di:           di,
		tokenService: tokenService,
		cacheService: cacheService,
	}, nil
}

func (s *sessionService) CreateSession(ctx context.Context, userID uuid.UUID, restaurantID *uuid.UUID, role models.Role) (*models.Session, error) {
	sessionID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	token, err := s.tokenService.CreateToken(userID, sessionID)
	if err != nil {
		return nil, err
	}

	session := &models.Session{
		UserID:       userID,
		SessionID:    sessionID,
		RestaurantID: restaurantID,
		Role:         role,
		Token:        token,
		CreatedAt:    time.Now().Unix(),
	}

	ttl := time.Duration(config.Env.Cache.SessionExp) * time.Hour

	if err := s.cacheService.Set(ctx, getSessionKey(sessionID), session, ttl); err != nil {
		return nil, err
	}

	if err := s.cacheService.AddToSet(ctx, getUserSessionsKey(userID), sessionID.String(), ttl); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *sessionService) GetSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	sessionID, err := s.tokenService.ExtractSessionID(token)
	if err != nil {
		return nil, models.ErrSessionNotFound
	}

	session, err := s.getSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.Token != token {
		return nil, models.ErrSessionNotFound
	}

	return session, nil
}

func (s *sessionService) GetSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Session, error) {
	var sessionIDs []string
	if err := s.cacheService.GetSetMembers(ctx, getUserSessionsKey(userID), &sessionIDs); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, models.ErrSessionNotFound
		}

		return nil, err
	}

	var activeSessions []*models.Session
	for _, sessionIDStr := range sessionIDs {
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			continue
		}

		session, err := s.getSession(ctx, sessionID)
		if err != nil {
			continue
		}

		if session != nil {
			_ = s.cacheService.RemoveFromSet(ctx, getUserSessionsKey(userID), sessionIDStr)
			continue
		}

		activeSessions = append(activeSessions, session)
	}

	return activeSessions, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.getSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session == nil {
		return models.ErrSessionNotFound
	}

	if err := s.cacheService.Delete(ctx, getSessionKey(sessionID)); err != nil {
		return err
	}

	return s.cacheService.RemoveFromSet(ctx, getUserSessionsKey(session.UserID), sessionID.String())
}

func (s *sessionService) DeleteAllSessions(ctx context.Context, userID uuid.UUID) error {
	var sessionIDs []string
	if err := s.cacheService.GetSetMembers(ctx, getUserSessionsKey(userID), &sessionIDs); err != nil {
		return err
	}

	if sessionIDs == nil {
		return models.ErrSessionNotFound
	}

	for _, sessionIDStr := range sessionIDs {
		sessionID, err := uuid.Parse(sessionIDStr)
		if err != nil {
			continue
		}
		_ = s.DeleteSession(ctx, sessionID)
	}

	return s.cacheService.Delete(ctx, getUserSessionsKey(userID))
}

func (s *sessionService) getSession(ctx context.Context, sessionID uuid.UUID) (*models.Session, error) {
	var session *models.Session
	if err := s.cacheService.Get(ctx, getSessionKey(sessionID), &session); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return nil, nil
		}

		return nil, err
	}

	return session, nil
}

func getSessionKey(sessionID uuid.UUID) string {
	return "session:" + sessionID.String()
}

func getUserSessionsKey(userID uuid.UUID) string {
	return "user_sessions:" + userID.String()
}
