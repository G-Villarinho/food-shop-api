package repositories

import (
	"context"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product models.Product) error
	GetProductByID(ctx context.Context, ID uuid.UUID) (*models.Product, error)
	GetProductsByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Product, error)
	GetProductsByIDsAndRestaurantID(ctx context.Context, productIDs []uuid.UUID, restaurantID uuid.UUID) ([]models.Product, error)
}

type productRepository struct {
	di *internal.Di
	DB *gorm.DB
}

func NewProductRepository(di *internal.Di) (ProductRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	return &productRepository{
		di: di,
		DB: db,
	}, nil
}

func (p *productRepository) CreateProduct(ctx context.Context, product models.Product) error {
	if err := p.DB.WithContext(ctx).Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func (p *productRepository) GetProductByID(ctx context.Context, ID uuid.UUID) (*models.Product, error) {
	var product *models.Product
	if err := p.DB.WithContext(ctx).Where("Id = ?", ID).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (p *productRepository) GetProductsByIDsAndRestaurantID(ctx context.Context, productIDs []uuid.UUID, restaurantID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	if err := p.DB.WithContext(ctx).Where("RestaurantID = ? AND ID IN (?)", restaurantID, productIDs).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (p *productRepository) GetProductsByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	if err := p.DB.WithContext(ctx).Where("RestaurantID = ?", restaurantID).Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
