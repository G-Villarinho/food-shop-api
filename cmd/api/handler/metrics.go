package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/food-shop-api/cmd/api/responses"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/services"
	"github.com/labstack/echo/v4"
)

type MetricsHandler interface {
	GetMonthlyMetrics(ctx echo.Context) error
}

type metricsHandler struct {
	di             *internal.Di
	metricsService services.MetricsService
}

func NewMetricsHandler(di *internal.Di) (MetricsHandler, error) {
	metricsService, err := internal.Invoke[services.MetricsService](di)
	if err != nil {
		return nil, err
	}

	return &metricsHandler{
		di:             di,
		metricsService: metricsService,
	}, nil
}

func (m *metricsHandler) GetMonthlyMetrics(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "metrics"),
		slog.String("func", "GetMonthlyMetrics"),
	)

	response, err := m.metricsService.GetMonthOrdersAmount(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrorOrderNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, 404, "not_found", "Nenhum pedido foi encontrado para a criação de métricas para o seu restaurante")
		}

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, 404, "not_found", "Restaurante não encontrado")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
