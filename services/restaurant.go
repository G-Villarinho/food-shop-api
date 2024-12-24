package services

import (
	"context"
	"fmt"

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
	orderService         OrderService
	userService          UserService
	restaurantRepository repositories.RestaurantRepository
}

func NewRestaurantService(di *internal.Di) (RestaurantService, error) {
	orderService, err := internal.Invoke[OrderService](di)
	if err != nil {
		return nil, err
	}

	userService, err := internal.Invoke[UserService](di)
	if err != nil {
		return nil, err
	}

	restaurantRepository, err := internal.Invoke[repositories.RestaurantRepository](di)
	if err != nil {
		return nil, err
	}

	return &restaurantService{
		di:                   di,
		orderService:         orderService,
		userService:          userService,
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
	custommerID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return models.ErrUserNotFoundInContext
	}

	restaurant, err := r.restaurantRepository.GetRestaurantByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("get restaurant by id: %w", err)
	}

	if restaurant == nil {
		return models.ErrRestaurantNotFound
	}

	if err := r.orderService.CreateOrder(ctx, custommerID, restaurantID, payload); err != nil {
		return err
	}

	return nil
}
