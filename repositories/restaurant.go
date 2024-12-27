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

type RestaurantRepository interface {
	CreateRestaurant(ctx context.Context, restaurant models.Restaurant) error
	GetRestaurantByID(ctx context.Context, ID uuid.UUID) (*models.Restaurant, error)
	GetRestaurantIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error)
	GetRestaurantByUserID(ctx context.Context, userID uuid.UUID) (*models.Restaurant, error)
}

type restaurantRepository struct {
	di *internal.Di
	DB *gorm.DB
}

func NewRestaurantRepository(di *internal.Di) (RestaurantRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	return &restaurantRepository{
		di: di,
		DB: db,
	}, nil
}

func (r *restaurantRepository) CreateRestaurant(ctx context.Context, restaurant models.Restaurant) error {
	if err := r.DB.Create(&restaurant).Error; err != nil {
		return err
	}

	return nil
}

func (r *restaurantRepository) GetRestaurantByID(ctx context.Context, ID uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.DB.Where("ID = ?", ID).First(&restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &restaurant, nil
}

func (r *restaurantRepository) GetRestaurantIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	var restaurantIDStr string
	if err := r.DB.
		WithContext(ctx).
		Table("Restaurants").
		Select("Id").
		Where("ManagerID = ?", userID).
		Scan(&restaurantIDStr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return nil, fmt.Errorf("error to parse restaurant ID: %w", err)
	}
	return &restaurantID, nil
}

func (r *restaurantRepository) GetRestaurantByUserID(ctx context.Context, userID uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	if err := r.DB.
		WithContext(ctx).
		Where("ManagerID = ?", userID).
		First(&restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &restaurant, nil
}
