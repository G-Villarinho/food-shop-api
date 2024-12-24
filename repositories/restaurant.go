package repositories

import (
	"context"
	"errors"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	CreateRestaurant(ctx context.Context, restaurant models.Restaurant) error
	GetRestaurantByID(ctx context.Context, ID uuid.UUID) (*models.Restaurant, error)
	GetRestaurantIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error)
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
	var restaurantID uuid.UUID
	if err := r.DB.
		WithContext(ctx).
		Table("Restaurant").
		Select("Id").
		Where("ManagerID = ?", userID).
		Scan(&restaurantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &restaurantID, nil
}
