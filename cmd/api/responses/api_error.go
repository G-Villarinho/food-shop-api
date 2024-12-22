package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	StatusCode int    `json:"status"`
	Title      string `json:"title"`
	Details    string `json:"details"`
}

func InternalServerAPIErrorResponse(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, ErrorResponse{
		StatusCode: http.StatusInternalServerError,
		Title:      "Internal Server Error",
		Details:    "Something went wrong. Please try again later.",
	})
}

func AccessDeniedAPIErrorResponse(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Title:      "Access Denied",
		Details:    "You need to be logged in to access this resource.",
	})
}

func ForbiddenPermissionAPIErrorResponse(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, ErrorResponse{
		StatusCode: http.StatusForbidden,
		Title:      "Permission Denied",
		Details:    "You do not have permission to perform this action.",
	})
}

func CannotBindPayloadAPIErrorResponse(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Title:      "Unable to Process Request",
		Details:    "We encountered an issue while trying to process your request. The data you provided is not in the expected format.",
	}
	return ctx.JSON(http.StatusUnprocessableEntity, errorResponse)
}

func NewCustomValidationAPIErrorResponse(ctx echo.Context, statusCode int, title, details string) error {
	errorResponse := ErrorResponse{
		StatusCode: statusCode,
		Title:      title,
		Details:    details,
	}

	return ctx.JSON(statusCode, errorResponse)
}
