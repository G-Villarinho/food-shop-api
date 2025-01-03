package models

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound = errors.New("session not found in the database")
)

type Session struct {
	UserID       uuid.UUID  `json:"userId"`
	SessionID    uuid.UUID  `json:"sessionId"`
	RestaurantID *uuid.UUID `json:"restaurantID,omitempty"`
	Role         Role       `json:"role"`
	Token        string     `json:"token"`
	CreatedAt    int64      `json:"createdAt"`
}
