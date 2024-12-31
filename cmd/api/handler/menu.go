package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/food-shop-api/cmd/api/responses"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/services"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type MenuHandler interface {
	UpdateMenu(ctx echo.Context) error
}

type menuHandler struct {
	di          *internal.Di
	menuService services.MenuService
}

func NewMenuHandler(di *internal.Di) (MenuHandler, error) {
	menuService, err := internal.Invoke[services.MenuService](di)
	if err != nil {
		return nil, err
	}

	return &menuHandler{
		di:          di,
		menuService: menuService,
	}, nil
}

func (m *menuHandler) UpdateMenu(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "menu"),
		slog.String("func", "UpdateMenu"),
	)

	var payload models.UpdateMenuPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := m.menuService.UpdateMenu(ctx.Request().Context(), &payload); err != nil {
		log.Error("Error to update menu", slog.String("error", err.Error()))

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "O restaurante especificado não foi encontrado. Verifique o ID e tente novamente.")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusNoContent)
}
