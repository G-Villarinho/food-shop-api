package repositories

import (
	"context"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"gorm.io/gorm"
)

type EvaluationRepository interface {
	CreateEvaluation(ctx context.Context, evaluation models.Evaluation) error
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
