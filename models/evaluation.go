package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEvaluationNotFound                  = errors.New("evaluation not found in the database")
	ErrEvaluationDoesNotBelongToRestaurant = errors.New("evaluation does not belong to the restaurant")
)

type Evaluation struct {
	BaseModel
	CustommerID  uuid.UUID      `gorm:"column:CustommerID;type:char(36);not null"`
	RestaurantID uuid.UUID      `gorm:"column:RestaurantID;type:char(36);not null"`
	Custommer    User           `gorm:"foreignKey:CustommerID;references:ID;OnDelete:CASCADE"`
	Restaurant   Restaurant     `gorm:"foreignKey:RestaurantID;references:ID;OnDelete:CASCADE"`
	Rating       int            `gorm:"column:Rating;type:int;not null"`
	Comment      string         `gorm:"column:Comment;type:text;not null"`
	Answer       sql.NullString `gorm:"column:Answer;type:text;null;default:null"`
}

func (e *Evaluation) TableName() string {
	return "Evaluation"
}

type CreateEvaluationPayload struct {
	RestaurantID uuid.UUID `json:"restaurantId" validate:"required"`
	Comment      string    `json:"comment" validate:"required,max=500"`
	Rating       int       `json:"rating" validate:"required,gte=1,lte=5"`
}

type UpdateAnswerPayload struct {
	EvaluationID uuid.UUID `json:"evaluationId" validate:"required"`
	Answer       string    `json:"answer" validate:"required,max=500"`
}

type EvaluationPagination struct {
	Pagination
	CustomerName *string `json:"customerName"`
	Rating       *int    `json:"rating"`
}

type EvaluationResponse struct {
	ID            uuid.UUID `json:"id"`
	CustommerName string    `json:"custommerName"`
	Rating        int       `json:"rating"`
	Comment       string    `json:"comment"`
	Answer        *string   `json:"answer"`
	CreatedAt     string    `json:"createdAt"`
}

type EvaluationSummary struct {
	Rating int
	Total  int
}

type StarCount struct {
	Stars      int `json:"rating"`
	TotalStars int `json:"count"`
}

type EvaluationSummaryResponse struct {
	StarSummary []StarCount `json:"starSummary"`
	Average     float64     `json:"average"`
}

func (c *CreateEvaluationPayload) ToEvaluation(custommerID uuid.UUID) *Evaluation {
	ID, _ := uuid.NewV7()

	return &Evaluation{
		BaseModel: BaseModel{
			ID: ID,
		},
		CustommerID:  custommerID,
		RestaurantID: c.RestaurantID,
		Rating:       c.Rating,
		Comment:      c.Comment,
	}
}

func (e *Evaluation) ToEvaluationResponse() *EvaluationResponse {
	return &EvaluationResponse{
		ID:            e.ID,
		CustommerName: e.Custommer.FullName,
		Rating:        e.Rating,
		Comment:       e.Comment,
		Answer:        &e.Answer.String,
		CreatedAt:     e.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
