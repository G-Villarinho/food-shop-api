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
	return "Product"
}
