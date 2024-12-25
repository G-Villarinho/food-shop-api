package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/cmd/api/validation"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type AuthHandler interface {
	SignIn(ctx echo.Context) error
	VeryfyMagicLink(ctx echo.Context) error
	SignOut(ctx echo.Context) error
}

type authHandler struct {
	di          *internal.Di
	authService services.AuthService
}

func NewAuthHandler(di *internal.Di) (AuthHandler, error) {
	authService, err := internal.Invoke[services.AuthService](di)
	if err != nil {
		return nil, err
	}

	return &authHandler{
		di:          di,
		authService: authService,
	}, nil
}

func (a *authHandler) SignIn(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "auth"),
		slog.String("func", "SignIn"),
	)

	var payload models.SignInPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := validation.ValidateStruct(payload); err != nil {
		log.Warn("Error to validate JSON payload")
		return responses.NewValidationErrorResponse(ctx, err)
	}

	if err := a.authService.SignIn(ctx.Request().Context(), payload.Email); err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrUserNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, 404, "not_found", "User not found.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusOK)
}

func (a *authHandler) VeryfyMagicLink(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "auth"),
		slog.String("func", "VerifyMagicLink"),
	)
	code, err := uuid.Parse(ctx.QueryParam("code"))
	if err != nil {
		log.Warn("Invalid Magic Link code format")
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "invalid_request", "Invalid Magic Link code format.")
	}

	redirectURL := ctx.QueryParam("redirect")
	if redirectURL == "" {
		log.Warn("Redirect URL is missing")
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "invalid_request", "Redirect URL is required.")
	}

	if redirectURL != config.Env.RedirectURL {
		log.Warn("Redirect URL is invalid")
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "invalid_request", "Redirect URL is invalid.")
	}

	token, err := a.authService.VeryfyMagicLink(ctx.Request().Context(), code)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrMagicLinkNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Magic Link not found.")
		}

		if errors.Is(err, models.ErrUserNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "User not found.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	cookie := new(http.Cookie)
	cookie.Name = config.Env.CookieName
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)

	return ctx.Redirect(http.StatusFound, redirectURL)
}

func (a *authHandler) SignOut(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "auth"),
		slog.String("func", "SignOut"),
	)

	if err := a.authService.SignOut(ctx.Request().Context()); err != nil {
		log.Error(err.Error())
		return responses.InternalServerAPIErrorResponse(ctx)
	}

	cookie := new(http.Cookie)
	cookie.Name = config.Env.CookieName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)

	return ctx.NoContent(http.StatusOK)
}
