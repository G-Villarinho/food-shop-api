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

type EvaluationRepository interface {
	CreateEvaluation(ctx context.Context, evaluation models.Evaluation) error
	GetPaginatedEvaluationsByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.EvaluationPagination) (*models.PaginatedResponse[models.Evaluation], error)
}

type evaluationRepository struct {
	di *internal.Di
	DB *gorm.DB
}

func NewEvaluationRepository(di *internal.Di) (EvaluationRepository, error) {
	db, err := internal.Invoke[*gorm.DB](di)
	if err != nil {
		return nil, err
	}

	return &evaluationRepository{
		di: di,
		DB: db,
	}, nil
}

func (e *evaluationRepository) CreateEvaluation(ctx context.Context, evaluation models.Evaluation) error {
	if err := e.DB.WithContext(ctx).Create(&evaluation).Error; err != nil {
		return err
	}

	return nil
}

func (e *evaluationRepository) GetPaginatedEvaluationsByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.EvaluationPagination) (*models.PaginatedResponse[models.Evaluation], error) {
	query := e.DB.WithContext(ctx).
		Model(&models.Evaluation{}).
		Preload("Custommer").
		Where("RestaurantID = ?", restaurantID)

	if pagination.Rating != nil {
		query = query.Where("Evaluations.Rating = ?", *pagination.Rating)
	}

	if pagination.CustomerName != nil {
		query = query.Joins("JOIN User ON Users.Id = Evaluations.CustommerID").
			Where("Users.FullName LIKE ?", fmt.Sprintf("%%%s%%", *pagination.CustomerName))
	}

	evaluations, err := paginate[models.Evaluation](query, &pagination.Pagination, &models.Evaluation{})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return evaluations, nil
}
