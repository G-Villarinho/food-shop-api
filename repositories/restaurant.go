package repositories

import (
	"context"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	CreateRestaurant(ctx context.Context, restaurant models.Restaurant) error
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
