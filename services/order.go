package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, custommerID, restaurantID uuid.UUID, payload models.CreateOrderPayload) error
	GetPaginatedOrdersByRestaurantID(ctx context.Context, pagination *models.Pagination) (*models.PaginatedResponse[*models.OrderResponse], error)
}

type orderService struct {
	di                   *internal.Di
	orderItemService     OrderItemService
	orderRepository      repositories.OrderRepository
	productRepository    repositories.ProductRepository
	restaurantRepository repositories.RestaurantRepository
}

func NewOrderService(di *internal.Di) (OrderService, error) {
	orderItemService, err := internal.Invoke[OrderItemService](di)
	if err != nil {
		return nil, err
	}

	orderRepository, err := internal.Invoke[repositories.OrderRepository](di)
	if err != nil {
		return nil, err
	}

	productRepository, err := internal.Invoke[repositories.ProductRepository](di)
	if err != nil {
		return nil, err
	}

	restaurantRepository, err := internal.Invoke[repositories.RestaurantRepository](di)
	if err != nil {
		return nil, err
	}

	return &orderService{
		di:                   di,
		orderItemService:     orderItemService,
		orderRepository:      orderRepository,
		productRepository:    productRepository,
		restaurantRepository: restaurantRepository,
	}, nil
}

func (o *orderService) CreateOrder(ctx context.Context, custommerID, restaurantID uuid.UUID, payload models.CreateOrderPayload) error {
	var productsIDs []uuid.UUID
	for _, item := range payload.Items {
		productsIDs = append(productsIDs, item.ProductID)
	}

	products, err := o.productRepository.GetProductsByIDsAndRestaurantID(ctx, productsIDs, payload.RestaurantID)
	if err != nil {
		return fmt.Errorf("get products by ids and restaurant id: %w", err)
	}

	orderItemSummary, err := o.orderItemService.ValidateAndCalculateOrderItems(ctx, products, payload.Items)
	if err != nil {
		return models.ErrSomeProductsNotFound
	}

	order := models.NewOrder(custommerID, payload.RestaurantID, orderItemSummary.TotalInCents)
	if err := o.orderRepository.CreateOrderWithItems(ctx, order, orderItemSummary.OrderItems); err != nil {
		return fmt.Errorf("error to create order: %w", err)
	}

	return nil
}

func (o *orderService) GetPaginatedOrdersByRestaurantID(ctx context.Context, pagination *models.Pagination) (*models.PaginatedResponse[*models.OrderResponse], error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	restaurant, err := o.restaurantRepository.GetRestaurantByID(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("get restaurant by ID: %w", err)
	}

	if restaurant == nil {
		return nil, models.ErrRestaurantNotFound
	}

	paginatedOrders, err := o.orderRepository.GetPaginatedOrdersByRestaurantID(ctx, restaurantID, pagination)
	if err != nil {
		return nil, fmt.Errorf("get paginated orders by restaurant ID: %w", err)
	}

	if paginatedOrders == nil {
		return nil, nil
	}

	paginatedOrdersResponse := models.MapPaginatedResult(paginatedOrders, func(order models.Order) *models.OrderResponse {
		return order.ToOrderResponse()
	})

	return paginatedOrdersResponse, nil
}
