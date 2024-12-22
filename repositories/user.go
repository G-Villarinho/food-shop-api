package repositories

import (
	"context"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
}

type userRepository struct {
	di *internal.Di
	DB *gorm.DB
}

func NewUserRepository(di *internal.Di) (UserRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		di: di,
		DB: db,
	}, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) error {
	if err := u.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}
