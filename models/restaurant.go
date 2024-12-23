package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Restaurant struct {
	BaseModel
	Name        string         `gorm:"column:Name;type:varchar(255);not null"`
	Description sql.NullString `gorm:"column:Description;type:text;default:null"`
	ManagerID   uuid.UUID      `gorm:"column:ManagerID;type:char(36);not null"`
	Manager     User           `gorm:"foreignKey:ManagerID;references:ID;OnDelete:CASCADE"`
}

func (r *Restaurant) TableName() string {
	return "Restaurant"
}

type CreateRestaurantPayload struct {
	Manager        CreateUserPayload `json:"manager" validate:"required"`
	RestaurantName string            `json:"restaurantName" validate:"required,max=255"`
}

func (payload *CreateRestaurantPayload) ToRestaurant(managerID uuid.UUID) *Restaurant {
	ID, _ := uuid.NewV7()
	return &Restaurant{
		BaseModel: BaseModel{
			ID: ID,
		},
		Name:      payload.RestaurantName,
		ManagerID: managerID,
	}
}
