package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/G-Villarinho/level-up-api/cmd/api/responses"
	"github.com/G-Villarinho/level-up-api/cmd/api/validation"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

type EvaluationHandler interface {
	CreateEvaluation(ctx echo.Context) error
}

type evaluationHandler struct {
	di *internal.Di
	services.EvaluationService
}

func NewEvaluationHandler(di *internal.Di) (EvaluationHandler, error) {
	evaluationService, err := internal.Invoke[services.EvaluationService](di)
	if err != nil {
		return nil, err
	}

	return &evaluationHandler{
		di:                di,
		EvaluationService: evaluationService,
	}, nil
}

func (e *evaluationHandler) CreateEvaluation(ctx echo.Context) error {
	log := slog.With(
		slog.String("handler", "evaluation"),
		slog.String("func", "CreateEvaluation"),
	)

	var payload models.CreateEvaluationPayload
	if err := jsoniter.NewDecoder(ctx.Request().Body).Decode(&payload); err != nil {
		log.Warn("Error to decode JSON payload", slog.String("error", err.Error()))
		return responses.CannotBindPayloadAPIErrorResponse(ctx)
	}

	if err := validation.ValidateStruct(payload); err != nil {
		log.Warn("Error to validate JSON payload")
		return responses.NewValidationErrorResponse(ctx, err)
	}

	if err := e.EvaluationService.CreateEvaluation(ctx.Request().Context(), payload); err != nil {
		log.Error(err.Error())

		if errors.Is(err, models.ErrUserNotFoundInContext) {
			return responses.AccessDeniedAPIErrorResponse(ctx)
		}

		if errors.Is(err, models.ErrRestaurantNotFound) {
			return responses.NewCustomValidationAPIErrorResponse(ctx, http.StatusNotFound, "not_found", "Restaurante não encontrado para realizar a avaliação")
		}

		return responses.InternalServerAPIErrorResponse(ctx)
	}

	return ctx.NoContent(http.StatusCreated)
}
