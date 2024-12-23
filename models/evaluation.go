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
