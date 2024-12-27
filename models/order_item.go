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

	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
}

func (o *OrderItem) TableName() string {
	return "OrderItems"
}

type CreateOrderItemPayload struct {
	ProductID uuid.UUID `json:"productID" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type OrderItemSummary struct {
	OrderItems   []OrderItem
	TotalInCents int
}

func NewOrderItem(productID uuid.UUID, quantity, priceInCents int) *OrderItem {
	ID, _ := uuid.NewUUID()
	return &OrderItem{
		BaseModel: BaseModel{
			ID: ID,
		},
		ProductID:    productID,
		PriceInCents: priceInCents,
		Quantity:     quantity,
	}
}
