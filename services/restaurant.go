package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
)

type RestaurantService interface {
	CreateRestaurant(ctx context.Context, payload models.CreateRestaurantPayload) error
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
