package services

import (
	"context"
	"errors"
	"testing"

	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/mocks"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEvaluationService_CreateEvaluation(t *testing.T) {
	ctx := context.WithValue(context.Background(), internal.UserIDKey, uuid.New())

	t.Run("should create evaluation successfully", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}

		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateEvaluationPayload{
			RestaurantID: restaurantID,
			Rating:       5,
			Comment:      "Excellent food!",
		}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(&models.Restaurant{}, nil)

		evaluationRepository.On("CreateEvaluation", ctx, mock.MatchedBy(func(evaluation models.Evaluation) bool {
			return evaluation.CustommerID == custommerID &&
				evaluation.RestaurantID == restaurantID &&
				evaluation.Rating == payload.Rating &&
				evaluation.Comment == payload.Comment
		})).Return(nil)

		response, err := evaluationService.CreateEvaluation(ctx, payload)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, payload.Rating, response.Rating)
		assert.Equal(t, payload.Comment, response.Comment)

		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		evaluationRepository.AssertCalled(t, "CreateEvaluation", ctx, mock.AnythingOfType("models.Evaluation"))
	})

	t.Run("should return error when user ID is not in context", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}

		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
			restaurantRepository: restaurantRepository,
		}

		invalidCtx := context.Background()
		payload := models.CreateEvaluationPayload{
			RestaurantID: uuid.New(),
			Rating:       4,
			Comment:      "Great service!",
		}

		response, err := evaluationService.CreateEvaluation(invalidCtx, payload)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrUserNotFoundInContext)

		restaurantRepository.AssertNotCalled(t, "GetRestaurantByID", mock.Anything, mock.Anything)
		evaluationRepository.AssertNotCalled(t, "CreateEvaluation", mock.Anything, mock.Anything)
	})

	t.Run("should return error when restaurant is not found", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}

		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
			restaurantRepository: restaurantRepository,
		}

		restaurantID := uuid.New()
		payload := models.CreateEvaluationPayload{
			RestaurantID: restaurantID,
			Rating:       3,
			Comment:      "Average experience.",
		}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(nil, nil)

		response, err := evaluationService.CreateEvaluation(ctx, payload)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		evaluationRepository.AssertNotCalled(t, "CreateEvaluation", mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to create evaluation", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		restaurantRepository := &mocks.RestaurantRepository{}

		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
			restaurantRepository: restaurantRepository,
		}

		custommerID := ctx.Value(internal.UserIDKey).(uuid.UUID)
		restaurantID := uuid.New()
		payload := models.CreateEvaluationPayload{
			RestaurantID: restaurantID,
			Rating:       2,
			Comment:      "Not great.",
		}

		restaurantRepository.On("GetRestaurantByID", ctx, restaurantID).Return(&models.Restaurant{}, nil)

		evaluationRepository.On("CreateEvaluation", ctx, mock.MatchedBy(func(evaluation models.Evaluation) bool {
			return evaluation.CustommerID == custommerID &&
				evaluation.RestaurantID == restaurantID &&
				evaluation.Rating == payload.Rating &&
				evaluation.Comment == payload.Comment
		})).Return(errors.New("database error"))

		response, err := evaluationService.CreateEvaluation(ctx, payload)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "create evaluation")

		restaurantRepository.AssertCalled(t, "GetRestaurantByID", ctx, restaurantID)
		evaluationRepository.AssertCalled(t, "CreateEvaluation", ctx, mock.AnythingOfType("models.Evaluation"))
	})
}

func TestEvaluationService_GetPaginatedEvaluationsByRestaurantID(t *testing.T) {
	restaurantID := uuid.New()
	ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

	t.Run("should return paginated evaluations successfully", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)

		pagination := &models.EvaluationPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		evaluations := []models.Evaluation{
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				Custommer: models.User{FullName: "John Doe"},
				Rating:    5,
				Comment:   "Excellent!",
			},
			{
				BaseModel: models.BaseModel{ID: uuid.New()},
				Custommer: models.User{FullName: "Jane Doe"},
				Rating:    4,
				Comment:   "Very good!",
			},
		}

		paginatedResult := &models.PaginatedResponse[models.Evaluation]{
			Data:       evaluations,
			Total:      2,
			TotalPages: 1,
			Page:       1,
			Limit:      10,
		}

		// Mock do repositório
		evaluationRepository.On("GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination).Return(paginatedResult, nil)

		// Executar
		response, err := evaluationService.GetPaginatedEvaluationsByRestaurantID(ctx, pagination)

		// Verificar
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2, len(response.Data))
		assert.Equal(t, "John Doe", response.Data[0].CustommerName)
		assert.Equal(t, 5, response.Data[0].Rating)
		assert.Equal(t, "Excellent!", response.Data[0].Comment)

		evaluationRepository.AssertCalled(t, "GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		invalidCtx := context.Background()
		pagination := &models.EvaluationPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		// Executar
		response, err := evaluationService.GetPaginatedEvaluationsByRestaurantID(invalidCtx, pagination)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		evaluationRepository.AssertNotCalled(t, "GetPaginatedEvaluationsByRestaurantID", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
		pagination := &models.EvaluationPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		// Mock para erro no repositório
		evaluationRepository.On("GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination).Return(nil, errors.New("database error"))

		// Executar
		response, err := evaluationService.GetPaginatedEvaluationsByRestaurantID(ctx, pagination)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "get paginated evaluations by restaurant ID")

		evaluationRepository.AssertCalled(t, "GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination)
	})

	t.Run("should return nil when no evaluations are found", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
		pagination := &models.EvaluationPagination{
			Pagination: models.Pagination{
				Page:  1,
				Limit: 10,
			},
		}

		// Mock para nenhum resultado encontrado
		evaluationRepository.On("GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination).Return(nil, nil)

		// Executar
		response, err := evaluationService.GetPaginatedEvaluationsByRestaurantID(ctx, pagination)

		// Verificar
		assert.NoError(t, err)
		assert.Nil(t, response)

		evaluationRepository.AssertCalled(t, "GetPaginatedEvaluationsByRestaurantID", ctx, *restaurantID, pagination)
	})
}

func TestEvaluationService_UpdateAnswer(t *testing.T) {
	restaurantID := uuid.New()
	ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

	t.Run("should update answer successfully", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
		evaluationID := uuid.New()
		payload := models.UpdateAnswerPayload{
			EvaluationID: evaluationID,
			Answer:       "Thank you for your feedback!",
		}

		evaluation := &models.Evaluation{
			BaseModel:    models.BaseModel{ID: evaluationID},
			RestaurantID: *restaurantID,
			Rating:       5,
			Comment:      "Great food!",
		}

		// Mock para retornar a avaliação existente
		evaluationRepository.On("GetEvaluationByID", ctx, evaluationID).Return(evaluation, nil)

		// Mock para atualizar a resposta
		evaluationRepository.On("UpdateAnswer", ctx, evaluationID, payload.Answer).Return(nil)

		// Executar
		response, err := evaluationService.UpdateAnswer(ctx, payload)

		// Verificar
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, payload.Answer, *response.Answer)

		evaluationRepository.AssertCalled(t, "GetEvaluationByID", ctx, evaluationID)
		evaluationRepository.AssertCalled(t, "UpdateAnswer", ctx, evaluationID, payload.Answer)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		invalidCtx := context.Background()
		payload := models.UpdateAnswerPayload{
			EvaluationID: uuid.New(),
			Answer:       "Thank you for your feedback!",
		}

		// Executar
		response, err := evaluationService.UpdateAnswer(invalidCtx, payload)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		evaluationRepository.AssertNotCalled(t, "GetEvaluationByID", invalidCtx, mock.Anything)
		evaluationRepository.AssertNotCalled(t, "UpdateAnswer", invalidCtx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when evaluation is not found", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		evaluationID := uuid.New()
		payload := models.UpdateAnswerPayload{
			EvaluationID: evaluationID,
			Answer:       "Thank you for your feedback!",
		}

		// Mock para avaliação não encontrada
		evaluationRepository.On("GetEvaluationByID", ctx, evaluationID).Return(nil, nil)

		// Executar
		response, err := evaluationService.UpdateAnswer(ctx, payload)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrEvaluationNotFound)

		evaluationRepository.AssertCalled(t, "GetEvaluationByID", ctx, evaluationID)
		evaluationRepository.AssertNotCalled(t, "UpdateAnswer", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when evaluation does not belong to restaurant", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		otherRestaurantID := uuid.New()
		evaluationID := uuid.New()
		payload := models.UpdateAnswerPayload{
			EvaluationID: evaluationID,
			Answer:       "Thank you for your feedback!",
		}

		evaluation := &models.Evaluation{
			BaseModel:    models.BaseModel{ID: evaluationID},
			RestaurantID: otherRestaurantID,
			Rating:       5,
			Comment:      "Great food!",
		}

		// Mock para retornar uma avaliação de outro restaurante
		evaluationRepository.On("GetEvaluationByID", ctx, evaluationID).Return(evaluation, nil)

		// Executar
		response, err := evaluationService.UpdateAnswer(ctx, payload)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrEvaluationDoesNotBelongToRestaurant)

		evaluationRepository.AssertCalled(t, "GetEvaluationByID", ctx, evaluationID)
		evaluationRepository.AssertNotCalled(t, "UpdateAnswer", ctx, mock.Anything, mock.Anything)
	})

	t.Run("should return error when repository fails to update answer", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)
		evaluationID := uuid.New()
		payload := models.UpdateAnswerPayload{
			EvaluationID: evaluationID,
			Answer:       "Thank you for your feedback!",
		}

		evaluation := &models.Evaluation{
			BaseModel:    models.BaseModel{ID: evaluationID},
			RestaurantID: *restaurantID,
			Rating:       5,
			Comment:      "Great food!",
		}

		// Mock para retornar a avaliação existente
		evaluationRepository.On("GetEvaluationByID", ctx, evaluationID).Return(evaluation, nil)

		// Mock para falha ao atualizar a resposta
		evaluationRepository.On("UpdateAnswer", ctx, evaluationID, payload.Answer).Return(errors.New("database error"))

		// Executar
		response, err := evaluationService.UpdateAnswer(ctx, payload)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "update answer")

		evaluationRepository.AssertCalled(t, "GetEvaluationByID", ctx, evaluationID)
		evaluationRepository.AssertCalled(t, "UpdateAnswer", ctx, evaluationID, payload.Answer)
	})
}

func TestEvaluationService_GetEvaluationSumary(t *testing.T) {
	restaurantID := uuid.New()
	ctx := context.WithValue(context.Background(), internal.RestaurantIDKey, &restaurantID)

	t.Run("should return evaluation summary successfully", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)

		summaries := []models.EvaluationSummary{
			{Rating: 5, Total: 10},
			{Rating: 4, Total: 5},
			{Rating: 3, Total: 2},
		}

		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, *restaurantID).Return(summaries, nil)

		response, err := evaluationService.GetEvaluationSumary(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 17, response.TotalCount)
		assert.Equal(t, 5, len(response.StarSummary))

		expectedAverage := (float64(5*10+4*5+3*2) / float64(17))
		assert.InDelta(t, expectedAverage, response.Average, 0.01)

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, *restaurantID)
	})

	t.Run("should handle all evaluations with the same rating", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		summaries := []models.EvaluationSummary{
			{Rating: 5, Total: 10},
		}

		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, restaurantID).Return(summaries, nil)

		response, err := evaluationService.GetEvaluationSumary(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 10, response.TotalCount)
		assert.Equal(t, 5.0, response.Average)
		assert.Equal(t, 5, len(response.StarSummary))
		assert.Equal(t, 100.0, response.StarSummary[0].Percentage)

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, restaurantID)
	})

	t.Run("should handle all evaluations with one-star ratings", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		summaries := []models.EvaluationSummary{
			{Rating: 1, Total: 15},
		}

		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, restaurantID).Return(summaries, nil)

		response, err := evaluationService.GetEvaluationSumary(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 15, response.TotalCount)
		assert.Equal(t, 1.0, response.Average)
		assert.Equal(t, 5, len(response.StarSummary))
		assert.Equal(t, 100.0, response.StarSummary[4].Percentage)

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, restaurantID)
	})

	t.Run("should handle ratings with some missing star categories", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		summaries := []models.EvaluationSummary{
			{Rating: 5, Total: 8},
			{Rating: 3, Total: 4},
		}

		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, restaurantID).Return(summaries, nil)

		response, err := evaluationService.GetEvaluationSumary(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 12, response.TotalCount)
		assert.InDelta(t, 4.33, response.Average, 0.01)
		assert.Equal(t, 5, len(response.StarSummary))
		assert.Equal(t, 0.0, response.StarSummary[3].Percentage) // Checar estrelas ausentes

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, restaurantID)
	})

	t.Run("should handle a large number of evaluations", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		summaries := []models.EvaluationSummary{
			{Rating: 5, Total: 1000},
			{Rating: 4, Total: 800},
			{Rating: 3, Total: 500},
			{Rating: 2, Total: 200},
			{Rating: 1, Total: 100},
		}

		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, restaurantID).Return(summaries, nil)

		response, err := evaluationService.GetEvaluationSumary(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 2600, response.TotalCount)

		// Verificar média correta
		expectedAverage := float64(10200) / float64(2600) // 3.923076923076923
		assert.InDelta(t, expectedAverage, response.Average, 0.01)

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, restaurantID)
	})

	t.Run("should return error when restaurant ID is not in context", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		invalidCtx := context.Background()

		// Executar
		response, err := evaluationService.GetEvaluationSumary(invalidCtx)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrRestaurantNotFound)

		evaluationRepository.AssertNotCalled(t, "GetEvaluationSumaryByRestaurantID", invalidCtx, mock.Anything)
	})

	t.Run("should return error when no summaries are found", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)

		// Mock para retornar nil como lista de summaries
		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, *restaurantID).Return(nil, nil)

		// Executar
		response, err := evaluationService.GetEvaluationSumary(ctx)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.ErrorIs(t, err, models.ErrEvaluationNotFound)

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, *restaurantID)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		evaluationRepository := &mocks.EvaluationRepository{}
		evaluationService := &evaluationService{
			evaluationRepository: evaluationRepository,
		}

		restaurantID := ctx.Value(internal.RestaurantIDKey).(*uuid.UUID)

		// Mock para erro no repositório
		evaluationRepository.On("GetEvaluationSumaryByRestaurantID", ctx, *restaurantID).Return(nil, errors.New("database error"))

		// Executar
		response, err := evaluationService.GetEvaluationSumary(ctx)

		// Verificar
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Contains(t, err.Error(), "get evaluation summary")

		evaluationRepository.AssertCalled(t, "GetEvaluationSumaryByRestaurantID", ctx, *restaurantID)
	})
}
