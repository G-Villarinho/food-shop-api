package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrorOrderNotFound                  = errors.New("order not found in database")
	ErrOrderIsDelivered                 = errors.New("order is already delivered")
	ErrorOrderCannotBeCancelled         = errors.New("order cannot be cancelled unless status is 'pending' or 'processing'")
	ErrorOrderDoesNotBelongToRestaurant = errors.New("order does not belong to the specified restaurant")
	ErrOrderCannotBeApproved            = errors.New("order cannot be approved unless status is 'pending'")
	ErrOrderCannotBeDispatched          = errors.New("order cannot be dispatched unless status is 'processing'")
)

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
	Items []CreateOrderItemPayload `json:"items" validate:"required,dive,required"`
}

type OrderPagination struct {
	Pagination
	Status       *string `json:"status"`
	OrderID      *string `json:"orderId"`
	CustomerName *string `json:"customerName"`
}

type OrderResponse struct {
	ID            uuid.UUID   `json:"id"`
	CustommerName string      `json:"custommerName"`
	Status        OrderStatus `json:"status"`
	TotalInCents  int         `json:"totalInCents"`
	CreatedAt     string      `json:"createdAt"`
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

func (o *Order) ToOrderResponse() *OrderResponse {
	return &OrderResponse{
		ID:            o.ID,
		CustommerName: o.Custommer.FullName,
		Status:        o.Status,
		TotalInCents:  o.TotalInCents,
		CreatedAt:     o.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
