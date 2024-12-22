package responses

import (
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/validation"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	StatusCode int               `json:"status"`
	Title      string            `json:"title"`
	Details    string            `json:"details"`
	Errors     []ValidationError `json:"errors"`
}

func NewValidationErrorResponse(ctx echo.Context, validationErrors validation.ValidationErrors) error {
	errors := make([]ValidationError, 0, len(validationErrors))
	for field, message := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   field,
			Message: message,
		})
	}

	return ctx.JSON(http.StatusUnprocessableEntity, ValidationErrorResponse{
		StatusCode: http.StatusUnprocessableEntity,
		Title:      "Validation Error",
		Details:    "One or more fields are invalid.",
		Errors:     errors,
	})
}
