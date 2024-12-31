package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound      = errors.New("product not found in the database")
	ErrSomeProductsNotFound = errors.New("some products not found for this restaurant")
)

type Product struct {
	BaseModel
	Name         string         `gorm:"column:Name;type:varchar(255);not null"`
	Description  sql.NullString `gorm:"column:Description;type:varchar(400);default:null"`
	PriceInCents int            `gorm:"column:PriceInCents;type:int;not null"`
	RestaurantID uuid.UUID      `gorm:"column:RestaurantID;type:char(36);not null"`
	Restaurant   Restaurant     `gorm:"foreignKey:RestaurantID;references:ID;OnDelete:CASCADE"`
}

func (p *Product) TableName() string {
	return "Products"
}

type CreateOrUpdateProductPayload struct {
	Id          *uuid.UUID `json:"id"`
	Name        *string    `json:"name" validate:"required,min=1,max=255"`
	Description *string    `json:"description" validate:"min=1,max=400"`
	Price       *float32   `json:"price" validate:"required,min=0"`
}

type UpdateMenuPayload struct {
	Products          []CreateOrUpdateProductPayload `json:"products" validate:"required,dive"`
	DeletedProductIDs []uuid.UUID                    `json:"deletedProductIDs"`
}

type PopularProduct struct {
	Name  string
	Count int
}

type PopularProductResponse struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (p *PopularProduct) ToPopularProductResponse() *PopularProductResponse {
	return &PopularProductResponse{
		Name:  p.Name,
		Count: p.Count,
	}
}

func (coup *CreateOrUpdateProductPayload) ToProduct() *Product {
	ID, _ := uuid.NewV7()

	return &Product{
		BaseModel: BaseModel{
			ID: ID,
		},
		Name:         *coup.Name,
		Description:  sql.NullString{String: *coup.Description, Valid: coup.Description != nil},
		PriceInCents: int(*coup.Price * 100),
	}
}

func (p *Product) ApplyUpdatePayload(payload *CreateOrUpdateProductPayload) {
	if payload.Name != nil {
		p.Name = *payload.Name
	}

	if payload.Description != nil {
		p.Description = sql.NullString{String: *payload.Description, Valid: payload.Description != nil}
	}

	if payload.Price != nil {
		p.PriceInCents = int(*payload.Price * 100)
	}
}
