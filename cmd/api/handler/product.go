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

type ProductHandler interface {
	GetPopularProducts(ctx echo.Context) error
}

type productHandler struct {
	di             *internal.Di
	productService services.ProductService
}

func NewProductHandler(di *internal.Di) (ProductHandler, error) {
	productService, err := internal.Invoke[services.ProductService](di)
	if err != nil {
		return nil, err
	}

	return &productHandler{
		di:             di,
		productService: productService,
	}, nil
}

func (p *productHandler) GetPopularProducts(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "products"),
		slog.String("func", "GetPopularProducts"),
	)

	response, err := p.productService.GetPopularProducts(ctx.Request().Context())
	if err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrProductNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Nenhum produto popular foi encontrado")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.JSON(http.StatusOK, response)
}
