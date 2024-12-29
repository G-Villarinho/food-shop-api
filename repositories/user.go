package repositories

import (
	"context"
	"errors"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepository --output=../mocks --outpkg=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error)
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

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user *models.User
	if err := u.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepository) GetUserByID(ctx context.Context, ID uuid.UUID) (*models.User, error) {
	var user *models.User
	if err := u.DB.WithContext(ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
