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

type OrderHandler interface {
	CreateOrder(ctx echo.Context) error
}

type orderHandler struct {
	di           *internal.Di
	orderService services.OrderService
}

func NewOrderHandler(di *internal.Di) (OrderHandler, error) {
	orderService, err := internal.Invoke[services.OrderService](di)
	if err != nil {
		return nil, err
	}

	return &orderHandler{
		di:           di,
		orderService: orderService,
	}, nil
}

func (o *orderHandler) CreateOrder(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "restaurant"),
		slog.String("func", "CreateOrder"),
	)

	var payload models.CreateOrderPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := validation.ValidateStruct(payload); err != nil {
		log.Warn("Error to validate JSON payload")
		return responses.NewValidationErrorResponse(ctx, err)
	}

	if err := o.orderService.CreateOrder(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Restaurant not found")
		}

		if errors.Is(err, models.ErrSomeProductsNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "bad_request", "Some products not found")
		}

		if errors.Is(err, models.ErrUserNotFoundInContext) {
			return responses.AccessDeniedAPIErrorResponse(ctx)
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}