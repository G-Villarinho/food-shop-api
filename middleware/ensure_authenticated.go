package middleware

import (
	"context"
	"errors"
	"log/slog"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
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

			authHeader := ctx.Request().Header.Get("Authorization")

			if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
				return responses.AccessDeniedAPIErrorResponse(ctx)
			}

			response, err := sessionService.GetSessionByToken(ctx.Request().Context(), authHeader[7:])
			if err != nil {
				slog.Error(err.Error())

				if errors.Is(err, models.ErrSessionNotFound) {
					return responses.AccessDeniedAPIErrorResponse(ctx)
				}

				return responses.InternalServerAPIErrorResponse(ctx)
			}

			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.UserIDKey, response.UserID)))
			ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), internal.SessionIDKey, response.SessionID)))
			return next(ctx)
		}
	}
}
