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
	CreateEvaluation(ctx context.Context, payload models.CreateEvaluationPayload) (*models.EvaluationResponse, error)
	GetPaginatedEvaluationsByRestaurantID(ctx context.Context, pagination *models.EvaluationPagination) (*models.PaginatedResponse[*models.EvaluationResponse], error)
	UpdateAnswer(ctx context.Context, payload models.UpdateAnswerPayload) (*models.EvaluationResponse, error)
	GetEvaluationSumary(ctx context.Context) (*models.EvaluationSummaryResponse, error)
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

func (e *evaluationService) CreateEvaluation(ctx context.Context, payload models.CreateEvaluationPayload) (*models.EvaluationResponse, error) {
	custommerID, ok := ctx.Value(internal.UserIDKey).(uuid.UUID)
	if !ok {
		return nil, models.ErrUserNotFoundInContext
	}

	restaurant, err := e.restaurantRepository.GetRestaurantByID(ctx, payload.RestaurantID)
	if err != nil {
		return nil, fmt.Errorf("get restaurant by ID: %w", err)
	}

	if restaurant == nil {
		return nil, models.ErrRestaurantNotFound
	}

	evaluation := payload.ToEvaluation(custommerID)
	if err := e.evaluationRepository.CreateEvaluation(ctx, *evaluation); err != nil {
		return nil, fmt.Errorf("create evaluation: %w", err)
	}

	return evaluation.ToEvaluationResponse(), nil
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

func (e *evaluationService) UpdateAnswer(ctx context.Context, payload models.UpdateAnswerPayload) (*models.EvaluationResponse, error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	evaluation, err := e.evaluationRepository.GetEvaluationByID(ctx, payload.EvaluationID)
	if err != nil {
		return nil, fmt.Errorf("get evaluation by ID: %w", err)
	}

	if evaluation == nil {
		return nil, models.ErrEvaluationNotFound
	}

	if evaluation.RestaurantID != *restaurantID {
		return nil, models.ErrEvaluationDoesNotBelongToRestaurant
	}

	if err := e.evaluationRepository.UpdateAnswer(ctx, payload.EvaluationID, payload.Answer); err != nil {
		return nil, fmt.Errorf("update answer: %w", err)
	}

	evaluationResponse := evaluation.ToEvaluationResponse()
	evaluationResponse.Answer = &payload.Answer

	return evaluationResponse, nil
}

func (e *evaluationService) GetEvaluationSumary(ctx context.Context) (*models.EvaluationSummaryResponse, error) {
	restaurantID, ok := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
	if !ok {
		return nil, models.ErrRestaurantNotFound
	}

	summaries, err := e.evaluationRepository.GetEvaluationSumaryByRestaurantID(ctx, *restaurantID)
	if err != nil {
		return nil, fmt.Errorf("get evaluation summary: %w", err)
	}

	if summaries == nil {
		return nil, models.ErrEvaluationNotFound
	}

	return e.buildEvaluationSummary(summaries), nil
}

func (e *evaluationService) buildEvaluationSummary(summaries []models.EvaluationSummary) *models.EvaluationSummaryResponse {
	totalMap := make(map[int]int)
	totalStars := 0
	totalCount := 0

	for _, summary := range summaries {
		totalMap[summary.Rating] = summary.Total
		totalStars += summary.Rating * summary.Total
		totalCount += summary.Total
	}

	var starSummary []models.StarCount
	for i := 1; i <= 5; i++ {
		starSummary = append(starSummary, models.StarCount{
			Stars:      i,
			TotalStars: totalMap[i],
		})
	}

	var average float64
	if totalCount > 0 {
		average = float64(totalStars) / float64(totalCount)
	}

	return &models.EvaluationSummaryResponse{
		StarSummary: starSummary,
		Average:     average,
	}
}
