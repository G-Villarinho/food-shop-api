package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error
	GetPaginatedOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.Pagination) (*models.PaginatedResponse[models.Order], error)
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

func (o *orderRepository) GetPaginatedOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.Pagination) (*models.PaginatedResponse[models.Order], error) {
	query := o.DB.WithContext(ctx).Model(&models.Order{}).Preload("Custommer").Where("RestaurantID = ?", restaurantID)
	orders, err := paginate[models.Order](query, pagination, &models.Order{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return orders, nil
}
