package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/google/uuid"
)

type RestaurantService interface {
	CreateRestaurant(ctx context.Context, payload models.CreateRestaurantPayload) error
	CreateOrder(ctx context.Context, restaurantID uuid.UUID, payload models.CreateOrderPayload) error
}

type restaurantService struct {
	di                   *internal.Di
	userService          UserService
	orderRepository      repositories.OrderRepository
	productRepository    repositories.ProductRepository
	restaurantRepository repositories.RestaurantRepository
}

func NewRestaurantService(di *internal.Di) (RestaurantService, error) {
	userService, err := internal.Invoke[UserService](di)
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

	return &restaurantService{
		di:                   di,
		userService:          userService,
		orderRepository:      orderRepository,
		productRepository:    productRepository,
		restaurantRepository: restaurantRepository,
	}, nil
}

func (r *restaurantService) CreateRestaurant(ctx context.Context, payload models.CreateRestaurantPayload) error {
	userID, err := r.userService.CreateUser(ctx, payload.Manager, models.Manager)
	if err != nil {
		return err
	}

	if err := r.restaurantRepository.CreateRestaurant(ctx, *payload.ToRestaurant(userID)); err != nil {
		return fmt.Errorf("create restaurant: %w", err)
	}

	return nil
}

func (r *restaurantService) CreateOrder(ctx context.Context, restaurantID uuid.UUID, payload models.CreateOrderPayload) error {
	log := slog.With(
		slog.String("service", "restaurant"),
		slog.String("func", "CreateOrder"),
	)

	custommerID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return models.ErrUserNotFoundInContext
	}

	var productsIDs []uuid.UUID
	for _, item := range payload.Items {
		productsIDs = append(productsIDs, item.ProductID)
	}

	products, err := r.productRepository.GetProductsByIDsAndRestaurantID(ctx, productsIDs, restaurantID)
	if err != nil {
		return fmt.Errorf("get products by ids and restaurant id: %w", err)
	}

	productMap := make(map[uuid.UUID]models.Product)
	for _, product := range products {
		productMap[product.ID] = product
	}

	var orderItems []models.OrderItem
	totalInCents := 0

	for _, item := range payload.Items {
		product, exists := productMap[item.ProductID]
		if !exists {
			log.Warn("Some products not found", slog.String("productID", item.ProductID.String()))
			return models.ErrSomeProductsNotFound
		}

		subtotal := product.PriceInCents * item.Quantity
		totalInCents += subtotal

		orderItems = append(orderItems, *models.NewOrderItem(product.ID, item.Quantity, product.PriceInCents))
	}

	order := models.NewOrder(custommerID, restaurantID, totalInCents)
	if err := r.orderRepository.CreateOrderWithItems(ctx, order, orderItems); err != nil {
		return fmt.Errorf("error to create order: %w", err)
	}

	return nil
}
