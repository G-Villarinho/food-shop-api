package models

import "github.com/google/uuid"

type OrderItem struct {
	BaseModel
	OrderID      uuid.UUID `gorm:"column:OrderID;type:char(36);not null"`
	ProductID    uuid.UUID `gorm:"column:ProductID;type:char(36);not null"`
	Order        Order     `gorm:"foreignKey:OrderID;references:ID;OnDelete:CASCADE"`
	Product      Product   `gorm:"foreignKey:ProductID;references:ID;OnDelete:CASCADE"`
	Quantity     int       `gorm:"column:Quantity;type:int;not null"`
	PriceInCents int       `gorm:"column:PriceInCents;type:int;not null"`
}

func (o *OrderItem) TableName() string {
	return "OrderItem"
}
