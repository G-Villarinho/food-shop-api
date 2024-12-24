package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	GetOrders(ctx echo.Context) error
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

func (o *orderHandler) GetOrders(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "restaurant"),
		slog.String("func", "GetOrders"),
	)

	pagination, err := models.NewPagination(ctx.QueryParam("page"), ctx.QueryParam("limit"), ctx.QueryParam("sort"))
	if err != nil {
		log.Warn("Error to create pagination", slog.String("error", err.Error()))
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, "conflict", "Invalid pagination parameters")
	}

	restaurantID, err := uuid.Parse(ctx.Param("restaurantID"))
	if err != nil {
		log.Warn("Error to parse restaurantID", slog.String("error", err.Error()))
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusUnprocessableEntity, "conflict", "Invalid restaurant ID parameter")
	}

	response, err := o.orderService.GetPaginatedOrdersByRestaurantID(ctx.Request().Context(), restaurantID, pagination)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Restaurant not found")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
