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

type RestaurantHandler interface {
	CreateRestaurant(ctx echo.Context) error
}

type restaurantHandler struct {
	di                *internal.Di
	restaurantService services.RestaurantService
}

func NewRestaurantHandler(di *internal.Di) (RestaurantHandler, error) {
	restaurantService, err := internal.Invoke[services.RestaurantService](di)
	if err != nil {
		return nil, err
	}

	return &restaurantHandler{
		di:                di,
		restaurantService: restaurantService,
	}, nil
}

func (r *restaurantHandler) CreateRestaurant(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "restaurant"),
		slog.String("func", "CreateRestaurant"),
	)

	var payload models.CreateRestaurantPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := validation.ValidateStruct(payload); err != nil {
		log.Warn("Error to validate JSON payload")
		return responses.NewValidationErrorResponse(ctx, err)
	}

	if err := r.restaurantService.CreateRestaurant(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrEmailAlreadyExists) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusConflict, "conflict", "The email already registered. Please try again with a different email.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}
