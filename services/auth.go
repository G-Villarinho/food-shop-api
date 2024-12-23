package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/G-Villarinho/level-up-api/cache"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/google/uuid"
)

type AuthService interface {
	SignIn(ctx context.Context, payload models.SignInPayload) error
	VeryfyMagicLink(ctx context.Context, code uuid.UUID) (string, error)
}

type authService struct {
	di              *internal.Di
	cache           cache.CacheService
	tokenService    TokenService
	sessionService  SessionService
	userRespository repositories.UserRepository
}

func NewAuthService(di *internal.Di) (AuthService, error) {
	cacheService, err := internal.Invoke[cache.CacheService](di)
	if err != nil {
		return nil, err
	}

	tokenService, err := internal.Invoke[TokenService](di)
	if err != nil {
		return nil, err
	}

	sessionService, err := internal.Invoke[SessionService](di)
	if err != nil {
		return nil, err
	}

	userRepository, err := internal.Invoke[repositories.UserRepository](di)
	if err != nil {
		return nil, err
	}

	return &authService{
		di:              di,
		cache:           cacheService,
		tokenService:    tokenService,
		sessionService:  sessionService,
		userRespository: userRepository,
	}, nil
}

func (a *authService) SignIn(ctx context.Context, payload models.SignInPayload) error {
	user, err := a.userRespository.GetUserByEmail(ctx, payload.Email)
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
	fmt.Println(magicLink)

	if err := a.cache.Set(ctx, getMagicLinkKey(code), user.ID.String(), 15*time.Minute); err != nil {
		return fmt.Errorf("set magic link: %w", err)
	}

	return nil
}

func (a *authService) VeryfyMagicLink(ctx context.Context, code uuid.UUID) (string, error) {
	var userID uuid.UUID
	if err := a.cache.Get(ctx, getMagicLinkKey(code), &userID); err != nil {
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

	if err := a.cache.Delete(ctx, getMagicLinkKey(code)); err != nil {
		return "", fmt.Errorf("delete magic link: %w", err)
	}

	session, err := a.sessionService.CreateSession(ctx, user.ID)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return session.Token, nil
}

func getMagicLinkKey(code uuid.UUID) string {
	return fmt.Sprintf("magic-link:%s", code.String())
}
