package models

import "github.com/google/uuid"

type OrderStatus string

const (
	Pending    OrderStatus = "pending"
	Canceled   OrderStatus = "canceled"
	Processing OrderStatus = "processing"
	Delivering OrderStatus = "delivering"
	Delivered  OrderStatus = "delivered"
)

type Order struct {
	BaseModel
	CustommerID  uuid.UUID   `gorm:"column:CustommerID;type:char(36);not null"`
	RestaurantID uuid.UUID   `gorm:"column:RestaurantID;type:char(36);not null"`
	Custommer    User        `gorm:"foreignKey:CustommerID;references:ID;OnDelete:CASCADE"`
	Restaurant   Restaurant  `gorm:"foreignKey:RestaurantID;references:ID;OnDelete:CASCADE"`
	Status       OrderStatus `gorm:"column:Status;type:enum('pending', 'canceled', 'processing', 'delivering', 'delivered');default:'pending';not null;index"`
	TotalInCents int         `gorm:"column:TotalInCents;type:int;not null"`
}

func (o *Order) TableName() string {
	return "Order"
}

type CreateOrderPayload struct {
	RestaurantID uuid.UUID                `json:"restaurantId" validate:"required"`
	Items        []CreateOrderItemPayload `json:"items" validate:"required,dive,required"`
}

func NewOrder(custommerID, restaurantID uuid.UUID, totalInCents int) *Order {
	ID, _ := uuid.NewUUID()
	return &Order{
		BaseModel: BaseModel{
			ID: ID,
		},
		CustommerID:  custommerID,
		RestaurantID: restaurantID,
		TotalInCents: totalInCents,
	}
}