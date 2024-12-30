package services

import (
	"context"
	"testing"

	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestOrderItemService_ValidateAndCalculateOrderItems(t *testing.T) {
	t.Run("should calculate order items successfully", func(t *testing.T) {
		orderItemService := &orderItemService{}

		product1ID := uuid.New()
		product2ID := uuid.New()
		products := []models.Product{
			{BaseModel: models.BaseModel{ID: product1ID}, PriceInCents: 1000},
			{BaseModel: models.BaseModel{ID: product2ID}, PriceInCents: 2000},
		}

		items := []models.CreateOrderItemPayload{
			{ProductID: product1ID, Quantity: 2},
			{ProductID: product2ID, Quantity: 1},
		}

		response, err := orderItemService.ValidateAndCalculateOrderItems(context.Background(), products, items)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.OrderItems))
		assert.Equal(t, 4000, response.TotalInCents)
	})

	t.Run("should return error for invalid product ID", func(t *testing.T) {
		orderItemService := &orderItemService{}

		productID := uuid.New()
		products := []models.Product{
			{BaseModel: models.BaseModel{ID: productID}, PriceInCents: 1500},
		}

		items := []models.CreateOrderItemPayload{
			{ProductID: uuid.New(), Quantity: 1},
		}

		response, err := orderItemService.ValidateAndCalculateOrderItems(context.Background(), products, items)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "is not available")
	})

	t.Run("should handle an empty list of products", func(t *testing.T) {
		orderItemService := &orderItemService{}

		products := []models.Product{}

		items := []models.CreateOrderItemPayload{
			{ProductID: uuid.New(), Quantity: 1},
		}

		response, err := orderItemService.ValidateAndCalculateOrderItems(context.Background(), products, items)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "is not available")
	})

	t.Run("should handle an empty list of order items", func(t *testing.T) {
		orderItemService := &orderItemService{}

		products := []models.Product{
			{BaseModel: models.BaseModel{ID: uuid.New()}, PriceInCents: 2000},
		}

		items := []models.CreateOrderItemPayload{}

		response, err := orderItemService.ValidateAndCalculateOrderItems(context.Background(), products, items)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 0, response.TotalInCents)
		assert.Equal(t, 0, len(response.OrderItems))
	})
}
