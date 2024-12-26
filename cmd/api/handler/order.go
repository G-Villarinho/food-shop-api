package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/G-Villarinho/level-up-api/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler interface {
	GetOrders(ctx echo.Context) error
	CancelOrder(ctx echo.Context) error
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
		log.Error(err.Error())
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "invalid_pagination", "Invalid pagination paramenter")
	}

	orderPagination := &models.OrderPagination{
		Pagination:   *pagination,
		Status:       utils.GetQueryStringPointer(ctx.QueryParam("status")),
		OrderID:      utils.GetQueryStringPointer(ctx.QueryParam("orderId")),
		CustomerName: utils.GetQueryStringPointer(ctx.QueryParam("customerName")),
	}

	response, err := o.orderService.GetPaginatedOrdersByRestaurantID(ctx.Request().Context(), orderPagination)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Restaurant not found")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}

func (o *orderHandler) CancelOrder(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "restaurant"),
		slog.String("func", "CancelOrder"),
	)

	orderID, err := uuid.Parse(ctx.Param("orderId"))
	if err != nil {
		log.Error(err.Error())
		return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "invalid_order_id", "Pedido inválido")
	}

	err = o.orderService.CancelOrder(ctx.Request().Context(), orderID)
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Restaurante não encontrado")
		}

		if errors.Is(err, models.ErrorOrderNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Pedido não encontrado no nosso sistema")
		}

		if errors.Is(err, models.ErrOrderIsDelivered) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "order_is_delivered", "O pedido já foi entregue")
		}

		if errors.Is(err, models.ErrorOrderDoesNotBelongToRestaurant) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "order_does_not_belong_to_restaurant", "O pedido não pertence ao restaurante especificado")
		}

		if errors.Is(err, models.ErrorOrderCannotBeCancelled) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusBadRequest, "order_cannot_be_cancelled", "O pedido não pode ser cancelado após ele já ter sido entregue, enviado ou cancelado anteriormente.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}
