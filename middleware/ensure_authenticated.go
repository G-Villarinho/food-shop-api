package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/food-shop-api/cmd/api/responses"
	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/services"
	"github.com/labstack/echo/v4"
)

func EnsureAuthenticated(di *internal.Di) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sessionService, err := internal.Invoke[services.SessionService](di)
			if err != nil {
				slog.Error(err.Error())
				return responses.InternalServerAPIErrorResponse(ctx)
			}

			authToken, err := getAuthToken(ctx)
			if err != nil {
				slog.Error(err.Error())
				clearAuthToken(ctx)
				return responses.AccessDeniedAPIErrorResponse(ctx)
			}

			response, err := sessionService.GetSessionByToken(ctx.Request().Context(), authToken)
			if err != nil {
				slog.Error(err.Error())

				if errors.Is(err, models.ErrSessionNotFound) {
					clearAuthToken(ctx)
					return responses.AccessDeniedAPIErrorResponse(ctx)
				}

				return responses.InternalServerAPIErrorResponse(ctx)
			}

			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.UserIDKey, response.UserID)))
			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.SessionIDKey, response.SessionID)))
			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.RoleKey, response.Role)))

			if response.Role == models.Manager {
				ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.RestaurantIDKey, response.RestaurantID)))
			}

			return next(ctx)
		}
	}
}

func clearAuthToken(ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = config.Env.CookieName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)
}

func getAuthToken(ctx echo.Context) (string, error) {
	cookie, err := ctx.Cookie(config.Env.CookieName)
	if err != nil {
		if errors.Is(err, echo.ErrCookieNotFound) {
			return "", models.ErrSessionNotFound
		}
		return "", err
	}

	authToken := cookie.Value
	if authToken == "" {
		return "", models.ErrSessionNotFound
	}

	return authToken, nil
}
