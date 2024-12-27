package models

import "github.com/google/uuid"

type Evaluation struct {
	BaseModel
	CustommerID  uuid.UUID  `gorm:"column:CustommerID;type:char(36);not null"`
	RestaurantID uuid.UUID  `gorm:"column:RestaurantID;type:char(36);not null"`
	Custommer    User       `gorm:"foreignKey:CustommerID;references:ID;OnDelete:CASCADE"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;references:ID;OnDelete:CASCADE"`
	Rating       int        `gorm:"column:Rating;type:int;not null"`
	Comment      string     `gorm:"column:Comment;type:text;not null"`
}

func (e *Evaluation) TableName() string {
	return "Evaluation"
}

type CreateEvaluationPayload struct {
	RestaurantID uuid.UUID `json:"restaurantId" validate:"required"`
	Comment      string    `json:"comment" validate:"required"`
	Rating       int       `json:"rating" validate:"required,gte=1,lte=5"`
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
