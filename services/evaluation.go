package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/google/uuid"
)

type EvaluationService interface {
	CreateEvaluation(ctx context.Context, evaluation models.CreateEvaluationPayload) error
}

type evaluationService struct {
	di                   *internal.Di
	evaluationRepository repositories.EvaluationRepository
	restaurantRepository repositories.RestaurantRepository
}

func NewEvaluationService(di *internal.Di) (EvaluationService, error) {
	evaluationRepository, err := internal.Invoke[repositories.EvaluationRepository](di)
	if err != nil {
		return nil, err
	}

	restaurantRepository, err := internal.Invoke[repositories.RestaurantRepository](di)
	if err != nil {
		return nil, err
	}

	return &evaluationService{
		di:                   di,
		evaluationRepository: evaluationRepository,
		restaurantRepository: restaurantRepository,
	}, nil
}

func (e *evaluationService) CreateEvaluation(ctx context.Context, evaluation models.CreateEvaluationPayload) error {
	custommerID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return models.ErrUserNotFoundInContext
	}

	restaurant, err := e.restaurantRepository.GetRestaurantByID(ctx, evaluation.RestaurantID)
	if err != nil {
		return fmt.Errorf("get restaurant by ID: %w", err)
	}

	if restaurant == nil {
		return models.ErrRestaurantNotFound
	}

	if err := e.evaluationRepository.CreateEvaluation(ctx, *evaluation.ToEvaluation(custommerID)); err != nil {
		return fmt.Errorf("create evaluation: %w", err)
	}

	return nil
}
