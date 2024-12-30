package services

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_GetPopularProducts(t *testing.T) {
	restaurantID := uuid.New()
	ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

	t.Run("should return popular products successfully", func(t *testing.T) {
		productRepository := &mocks.ProductRepository{}
		productService := &productService{
			popularProductRepository: productRepository,
		}

		popularProducts := []models.PopularProduct{
			{Name: "Pizza", Count: 10},
			{Name: "Burger", Count: 8},
		}

		productRepository.On("GetPopularProducts", ctx, restaurantID, 5).Return(popularProducts, nil)

		products, err := productService.GetPopularProducts(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 2)
		assert.Equal(t, "Pizza", products[0].Name)
		assert.Equal(t, 10, products[0].Count)
		assert.Equal(t, "Burger", products[1].Name)
		assert.Equal(t, 8, products[1].Count)

		productRepository.AssertCalled(t, "GetPopularProducts", ctx, restaurantID, 5)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		productRepository := &mocks.ProductRepository{}
		productService := &productService{
			popularProductRepository: productRepository,
		}

		invalidCtx := context.Background()

		products, err := productService.GetPopularProducts(invalidCtx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		productRepository.AssertNotCalled(t, "GetPopularProducts", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository returns an error", func(t *testing.T) {
		productRepository := &mocks.ProductRepository{}
		productService := &productService{
			popularProductRepository: productRepository,
		}

		productRepository.On("GetPopularProducts", ctx, restaurantID, 5).Return(nil, errors.New("database error"))

		products, err := productService.GetPopularProducts(ctx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.Contains(t, err.Error(), "get popular products")

		productRepository.AssertCalled(t, "GetPopularProducts", ctx, restaurantID, 5)
	})

	t.Run("should return error when no products are found", func(t *testing.T) {
		productRepository := &mocks.ProductRepository{}
		productService := &productService{
			popularProductRepository: productRepository,
		}

		productRepository.On("GetPopularProducts", ctx, restaurantID, 5).Return(nil, nil)

		products, err := productService.GetPopularProducts(ctx)

		assert.Error(t, err)
		assert.Nil(t, products)
		assert.ErrorIs(t, err, models.ErrProductNotFound)

		productRepository.AssertCalled(t, "GetPopularProducts", ctx, restaurantID, 5)
	})
}
