package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/google/uuid"
)

type OrderItemService interface {
	ValidateAndCalculateOrderItems(ctx context.Context, products []models.Product, items []models.CreateOrderItemPayload) (*models.OrderItemSummary, error)
}

type orderItemService struct {
	di *internal.Di
}

func NewOrderItemService(di *internal.Di) (OrderItemService, error) {
	return &orderItemService{
		di: di,
	}, nil
}

func (o *orderItemService) ValidateAndCalculateOrderItems(ctx context.Context, products []models.Product, items []models.CreateOrderItemPayload) (*models.OrderItemSummary, error) {
	productMap := make(map[uuid.UUID]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	var orderItems []models.OrderItem
	totalInCents := 0

	for _, item := range items {
		product, exists := productMap[item.ProductID]
		if !exists {
			return nil, fmt.Errorf("product with ID %s is not available", item.ProductID)
		}

		subtotal := product.PriceInCents * item.Quantity
		totalInCents += subtotal

		orderItem := models.NewOrderItem(product.ID, item.Quantity, product.PriceInCents)
		orderItems = append(orderItems, *orderItem)
	}

	summary := &models.OrderItemSummary{
		OrderItems:   orderItems,
		TotalInCents: totalInCents,
	}

	return summary, nil
}
