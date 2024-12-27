package services

import (
	"context"
	"fmt"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/G-Villarinho/food-shop-api/repositories"
	"github.com/google/uuid"
)

type EvaluationService interface {
	CreateEvaluation(ctx context.Context, evaluation models.CreateEvaluationPayload) error
	GetPaginatedEvaluationsByRestaurantID(ctx context.Context, pagination *models.EvaluationPagination) (*models.PaginatedResponse[*models.EvaluationResponse], error)
	UpdateAnswer(ctx context.Context, payload models.UpdateAnswerPayload) error
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

func (e *evaluationService) GetPaginatedEvaluationsByRestaurantID(ctx context.Context, pagination *models.EvaluationPagination) (*models.PaginatedResponse[*models.EvaluationResponse], error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	paginatedEvaluations, err := e.evaluationRepository.GetPaginatedEvaluationsByRestaurantID(ctx, *restaurantID, pagination)
	if err != nil {
		return nil, fmt.Errorf("get paginated evaluations by restaurant ID: %w", err)
	}

	if paginatedEvaluations == nil {
		return nil, nil
	}

	paginatedEvaluationsResponse := models.MapPaginatedResult(paginatedEvaluations, func(evaluation models.Evaluation) *models.EvaluationResponse {
		return evaluation.ToEvaluationResponse()
	})

	return paginatedEvaluationsResponse, nil
}

func (e *evaluationService) UpdateAnswer(ctx context.Context, payload models.UpdateAnswerPayload) error {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return models.ErrRestaurantNotFound
	}

	evaluation, err := e.evaluationRepository.GetEvaluationByID(ctx, payload.EvaluationID)
	if err != nil {
		return fmt.Errorf("get evaluation by ID: %w", err)
	}

	if evaluation == nil {
		return models.ErrEvaluationNotFound
	}

	if evaluation.RestaurantID != *restaurantID {
		return models.ErrEvaluationDoesNotBelongToRestaurant
	}

	if err := e.evaluationRepository.UpdateAnswer(ctx, payload.EvaluationID, payload.Answer); err != nil {
		return fmt.Errorf("update answer: %w", err)
	}

	return nil
}
