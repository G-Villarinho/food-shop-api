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
	CreateProduct(ctx context.Context, payload *models.CreateOrUpdateProductPayload) (*models.Product, error)
	CreateProducts(ctx context.Context, payload []models.CreateOrUpdateProductPayload) error
	UpdateProducts(ctx context.Context, payload []models.CreateOrUpdateProductPayload) error
	DeleteProducts(ctx context.Context, productIDs []uuid.UUID, restaurantID uuid.UUID) error
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

func (p *productService) CreateProduct(ctx context.Context, payload *models.CreateOrUpdateProductPayload) (*models.Product, error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	product := payload.ToProduct()
	product.RestaurantID = *restaurantID

	if err := p.popularProductRepository.CreateProduct(ctx, *product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

func (p *productService) CreateProducts(ctx context.Context, payload []models.CreateOrUpdateProductPayload) error {
	var products []models.Product
	for _, product := range payload {
		products = append(products, *product.ToProduct())
	}

	if err := p.popularProductRepository.CreateRange(ctx, products); err != nil {
		return fmt.Errorf("create many products: %w", err)
	}

	return nil
}

func (p *productService) UpdateProducts(ctx context.Context, payload []models.CreateOrUpdateProductPayload) error {
	if len(payload) == 0 {
		return nil
	}

	var updatedProducts = make(map[uuid.UUID]models.CreateOrUpdateProductPayload)
	var productsIds []uuid.UUID
	for _, product := range payload {
		productsIds = append(productsIds, *product.Id)
		updatedProducts[*product.Id] = product
	}

	products, err := p.popularProductRepository.GetProductsByIds(ctx, productsIds)
	if err != nil {
		return fmt.Errorf("get products by ids: %w", err)
	}

	for _, product := range products {
		if updatedProduct, ok := updatedProducts[product.ID]; ok {
			product.ApplyUpdatePayload(&updatedProduct)
		}
	}

	if err := p.popularProductRepository.UpdateRange(ctx, products); err != nil {
		return fmt.Errorf("update many products: %w", err)
	}

	return nil
}

func (p *productService) DeleteProducts(ctx context.Context, productIDs []uuid.UUID, restaurantID uuid.UUID) error {
	if err := p.popularProductRepository.DeleteRange(ctx, productIDs, restaurantID); err != nil {
		return fmt.Errorf("delete many products: %w", err)
	}

	return nil
}
