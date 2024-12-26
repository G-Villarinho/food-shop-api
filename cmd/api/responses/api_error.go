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
		Details:    "Ocorreu um erro interno no servidor. Por favor, tente novamente mais tarde.",
	})
}

func AccessDeniedAPIErrorResponse(ctx echo.Context) error {
	return ctx.JSON(http.StatusUnauthorized, ErrorResponse{
		StatusCode: http.StatusUnauthorized,
		Title:      "Access Denied",
		Details:    "Você precisa estar autenticado para acessar este recurso.",
	})
}

func ForbiddenPermissionAPIErrorResponse(ctx echo.Context) error {
	return ctx.JSON(http.StatusForbidden, ErrorResponse{
		StatusCode: http.StatusForbidden,
		Title:      "Permission Denied",
		Details:    "Você não tem permissão para acessar este recurso.",
	})
}

func CannotBindPayloadAPIErrorResponse(ctx echo.Context) error {
	errorResponse := ErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Title:      "Unable to Process Request",
		Details:    "Não foi possível processar os dados enviados. Verifique se as informações estão no formato correto.",
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
