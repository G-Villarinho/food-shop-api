package repositories

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error
}

type orderRepository struct {
	di *internal.Di
	DB *gorm.DB
}

func NewOrderRepository(di *internal.Di) (OrderRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	return &orderRepository{
		di: di,
		DB: db,
	}, nil
}

func (o *orderRepository) CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	return o.DB.Transaction(func(tx *gorm.DB) error {
		// Criar o pedido
		if err := tx.WithContext(ctx).Create(order).Error; err != nil {
			return fmt.Errorf("error to create order: %w", err)
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.WithContext(ctx).Create(&items).Error; err != nil {
			return fmt.Errorf("error to create order items: %w", err)
		}

		return nil
	})
}
