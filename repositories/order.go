package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error
	GetPaginatedOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.OrderPagination) (*models.PaginatedResponse[models.Order], error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error
	GetOrderByID(ctx context.Context, orderID uuid.UUID, preload bool) (*models.Order, error)
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

func (o *orderRepository) GetPaginatedOrdersByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.OrderPagination) (*models.PaginatedResponse[models.Order], error) {
	query := o.DB.WithContext(ctx).
		Model(&models.Order{}).
		Preload("Custommer").
		Where("RestaurantID = ?", restaurantID)

	if pagination.Status != nil {
		query = query.Where("Orders.Status = ?", *pagination.Status)
	}

	if pagination.OrderID != nil {
		query = query.Where("Id LIKE ?", fmt.Sprintf("%%%s%%", *pagination.OrderID))
	}

	if pagination.CustomerName != nil {
		query = query.Joins("JOIN Users ON Users.Id = Orders.CustommerID").
			Where("Users.FullName LIKE ?", fmt.Sprintf("%%%s%%", *pagination.CustomerName))
	}

	orders, err := paginate[models.Order](query, &pagination.Pagination, &models.Order{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return orders, nil
}

func (o *orderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status models.OrderStatus) error {
	return o.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("Id = ?", orderID).
		Update("Status", status).
		Error
}

func (o *orderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID, preload bool) (*models.Order, error) {
	query := o.DB.WithContext(ctx).
		Model(&models.Order{}).
		Where("Id = ?", orderID)

	if preload {
		query = query.Preload("Custommer").
			Preload("Restaurant")
	}

	var order models.Order
	if err := query.First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &order, nil
}
