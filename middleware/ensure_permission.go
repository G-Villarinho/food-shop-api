package middleware

import (
	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/labstack/echo/v4"
)

func EnsurePermission(requiredPermission models.Permission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			role, ok := ctx.Request().Context().Value(internal.RoleKey).(models.Role)
			if !ok || role == "" {
				return responses.AccessDeniedAPIErrorResponse(ctx)
			}

			if hasPermission := models.CheckPermission(role, requiredPermission); !hasPermission {
				return responses.ForbiddenPermissionAPIErrorResponse(ctx)
			}

			return next(ctx)
		}
	}
}
