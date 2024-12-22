package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors map[string]string

func ValidateStruct(s any) ValidationErrors {
	validate := validator.New()

	if err := SetupCustomValidations(validate); err != nil {
		return ValidationErrors{"validation_setup": "error to set up custom validations"}
	}

	validationErrors := make(ValidationErrors)
	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			validationErrors[fieldName] = getErrorMessage(err)
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return validationErrors
}

func getErrorMessage(err validator.FieldError) string {
	if msg, exists := ValidationMessages[err.Tag()]; exists {
		return msg
	}
	return "Invalid value"
}
