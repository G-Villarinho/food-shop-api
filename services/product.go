package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/repositories"
	"github.com/google/uuid"
)

//go:generate mockery --name=ProductService --output=../mocks --outpkg=mocks
type ProductService interface {
	GetPopularProducts(ctx context.Context) ([]models.PopularProductResponse, error)
}

type productService struct {
	di                       *internal.Di
	popularProductRepository repositories.ProductRepository
}

func NewProductService(di *internal.Di) (ProductService, error) {
	popularProductRepository, err := internal.Invoke[repositories.ProductRepository](di)
	if err != nil {
		return nil, err
	}

	return &productService{
		di:                       di,
		popularProductRepository: popularProductRepository,
	}, nil
}

func (p *productService) GetPopularProducts(ctx context.Context) ([]models.PopularProductResponse, error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	const popularProductsLimit = 5

	popularProducts, err := p.popularProductRepository.GetPopularProducts(ctx, *restaurantID, popularProductsLimit)
	if err != nil {
		return nil, fmt.Errorf("get popular products: %w", err)
	}

	if popularProducts == nil {
		return nil, models.ErrProductNotFound
	}

	var popularProductsResponse []models.PopularProductResponse
	for _, product := range popularProducts {
		popularProductsResponse = append(popularProductsResponse, *product.ToPopularProductResponse())
	}

	return popularProductsResponse, nil
}
