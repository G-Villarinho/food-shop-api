package services

import (
	"context"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
)

type MenuService interface {
	UpdateMenu(ctx context.Context, payload *models.UpdateMenuPayload) error
}

type menuService struct {
	di             *internal.Di
	productService ProductService
}

func NewMenuService(di *internal.Di) (MenuService, error) {
	productService, err := internal.Invoke[ProductService](di)
	if err != nil {
		return nil, err
	}

	return &menuService{
		di:             di,
		productService: productService,
	}, nil
}

func (m *menuService) UpdateMenu(ctx context.Context, payload *models.UpdateMenuPayload) error {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return models.ErrRestaurantNotFound
	}

	if len(payload.DeletedProductIDs) > 0 {
		if err := m.productService.DeleteMany(ctx, payload.DeletedProductIDs, *restaurantID); err != nil {
			return err
		}
	}

	var updatedProducts []*models.CreateOrUpdateProductPayload
	for _, product := range payload.Products {
		if product.Id != nil {
			updatedProducts = append(updatedProducts, &product)
		}
	}

	for _, product := range updatedProducts {
		if _, err := m.productService.UpdateProduct(ctx, product); err != nil {
			return err
		}
	}

	var newProducts []*models.CreateOrUpdateProductPayload
	for _, product := range payload.Products {
		if product.Id == nil {
			newProducts = append(newProducts, &product)
		}
	}

	for _, product := range newProducts {
		if _, err := m.productService.CreateProduct(ctx, product); err != nil {
			return err
		}
	}

	return nil
}
