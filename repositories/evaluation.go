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

//go:generate mockery --name=EvaluationRepository --output=../mocks --outpkg=mocks
type EvaluationRepository interface {
	CreateEvaluation(ctx context.Context, evaluation models.Evaluation) error
	GetPaginatedEvaluationsByRestaurantID(ctx context.Context, restaurantID uuid.UUID, pagination *models.EvaluationPagination) (*models.PaginatedResponse[models.Evaluation], error)
	UpdateAnswer(ctx context.Context, evaluationID uuid.UUID, answer string) error
	GetEvaluationByID(ctx context.Context, evaluationID uuid.UUID) (*models.Evaluation, error)
	GetEvaluationSumaryByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.EvaluationSummary, error)
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

func (e *evaluationRepository) UpdateAnswer(ctx context.Context, evaluationID uuid.UUID, answer string) error {
	if err := e.DB.WithContext(ctx).
		Model(&models.Evaluation{}).
		Where("Id = ?", evaluationID).
		Update("Answer", answer).Error; err != nil {
		return err
	}

	return nil
}

func (e *evaluationRepository) GetEvaluationByID(ctx context.Context, evaluationID uuid.UUID) (*models.Evaluation, error) {
	var evaluation models.Evaluation
	if err := e.DB.WithContext(ctx).
		Where("Id = ?", evaluationID).
		First(&evaluation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &evaluation, nil
}

func (e *evaluationRepository) GetEvaluationSumaryByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.EvaluationSummary, error) {
	var evaluationSummary []models.EvaluationSummary

	if err := e.DB.WithContext(ctx).
		Model(&models.Evaluation{}).
		Select("Rating, COUNT(*) as Total").
		Where("RestaurantID = ?", restaurantID).
		Group("Rating").
		Order("Rating").
		Scan(&evaluationSummary).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return evaluationSummary, nil
}
