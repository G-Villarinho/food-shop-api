package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/G-Villarinho/level-up-api/cache"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/G-Villarinho/level-up-api/services/email"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type AuthService interface {
	SignIn(ctx context.Context, email string) error
	VeryfyMagicLink(ctx context.Context, code uuid.UUID) (string, error)
	SignOut(ctx context.Context) error
}

type authService struct {
	di                   *internal.Di
	emailFactory         email.EmailFactory
	cacheService         cache.CacheService
	queueService         QueueService
	sessionService       SessionService
	restaurantRepository repositories.RestaurantRepository
	userRespository      repositories.UserRepository
}

func NewAuthService(di *internal.Di) (AuthService, error) {
	cacheService, err := internal.Invoke[cache.CacheService](di)
	if err != nil {
		return nil, err
	}

	queueService, err := internal.Invoke[QueueService](di)
	if err != nil {
		return nil, err
	}

	sessionService, err := internal.Invoke[SessionService](di)
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

	return &authService{
		di:                   di,
		emailFactory:         *email.NewEmailTaskFactory(),
		cacheService:         cacheService,
		queueService:         queueService,
		sessionService:       sessionService,
		restaurantRepository: restaurantRepository,
		userRespository:      userRepository,
	}, nil
}

func (a *authService) SignIn(ctx context.Context, email string) error {
	log := slog.With(
		slog.String("service", "auth"),
		slog.String("func", "SignIn"),
	)

	user, err := a.userRespository.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user == nil {
		return models.ErrUserNotFound
	}

	code, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate code: %w", err)
	}

	magicLink := fmt.Sprintf("%s/auth/link?code=%s&redirect=%s", config.Env.APIBaseURL, code.String(), config.Env.RedirectURL)

	if err := a.cacheService.Set(ctx, getMagicLinkKey(code), user.ID.String(), 15*time.Minute); err != nil {
		return fmt.Errorf("set magic link: %w", err)
	}

	go func() {
		message, err := jsoniter.Marshal(a.emailFactory.CreateSignInMagicLinkEmail(user.Email, user.FullName, magicLink))
		if err != nil {
			log.Error("marshal email task", slog.String("error", err.Error()))
			return
		}

		if err := a.queueService.Publish(QueueSendEmail, message); err != nil {
			log.Error("publish email task", slog.String("error", err.Error()))
			return
		}
	}()

	return nil
}

func (a *authService) VeryfyMagicLink(ctx context.Context, code uuid.UUID) (string, error) {
	var userID uuid.UUID
	if err := a.cacheService.Get(ctx, getMagicLinkKey(code), &userID); err != nil {
		if errors.Is(err, cache.ErrCacheMiss) {
			return "", models.ErrMagicLinkNotFound
		}
		return "", fmt.Errorf("get magic link: %w", err)
	}

	user, err := a.userRespository.GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get user by id: %w", err)
	}

	if user == nil {
		return "", models.ErrUserNotFound
	}

	if err := a.cacheService.Delete(ctx, getMagicLinkKey(code)); err != nil {
		return "", fmt.Errorf("delete magic link: %w", err)
	}

	var restaurantID *uuid.UUID
	if user.Role == models.Manager {
		restaurantID, err = a.restaurantRepository.GetRestaurantIDByUserID(ctx, user.ID)
		if err != nil {
			return "", fmt.Errorf("get restaurant id by user id: %w", err)
		}
	}

	session, err := a.sessionService.CreateSession(ctx, user.ID, restaurantID, user.Role)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return session.Token, nil
}

func (a *authService) SignOut(ctx context.Context) error {
	sessionId, ok := ctx.Value(internal.SessionIDKey).(uuid.UUID)
	if !ok {
		return models.ErrSessionNotFound
	}

	if err := a.sessionService.DeleteSession(ctx, sessionId); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}

func getMagicLinkKey(code uuid.UUID) string {
	return fmt.Sprintf("magic-link:%s", code.String())
}
