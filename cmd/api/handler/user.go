package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/cmd/api/validation"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	CreateUser(ctx echo.Context) error
	GetUser(ctx echo.Context) error
}

type userHandler struct {
	di          *internal.Di
	userService services.UserService
}

func NewUserHandler(di *internal.Di) (UserHandler, error) {
	userService, err := internal.Invoke[services.UserService](di)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		di:          di,
		userService: userService,
	}, nil
}

func (u *userHandler) CreateUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "GetUser"),
	)

	var payload models.CreateUserPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := validation.ValidateStruct(payload); err != nil {
		log.Warn("Error to validate JSON payload")
		return responses.NewValidationErrorResponse(ctx, err)
	}

	if err := u.userService.CreateUser(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrEmailAlreadyExists) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, "conflict", "The email already registered. Please try again with a different email.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (u *userHandler) GetUser(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "user"),
		slog.String("func", "GetUser"),
	)

	response, err := u.userService.GetUser(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrUserNotFound) {
			return responses.AccessDeniedAPIErrorResponse(ctx)
		}

		if errors.Is(err, models.ErrUserNotFoundInContext) {
			return responses.AccessDeniedAPIErrorResponse(ctx)
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
