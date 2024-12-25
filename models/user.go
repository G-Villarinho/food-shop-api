package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound          = errors.New("user not found in the database")
	ErrEmailAlreadyExists    = errors.New("email already exists in the database")
	ErrUserNotFoundInContext = errors.New("user not found in the context")
)

type Status string
type Role string

const (
	Active  Status = "active"
	Blocked Status = "blocked"
)

const (
	Manager  Role = "manager"
	Customer Role = "customer"
)

type User struct {
	BaseModel
	FullName string         `gorm:"column:FullName;type:varchar(255);not null"`
	Email    string         `gorm:"column:Email;type:varchar(255);not null;unique"`
	Status   Status         `gorm:"column:Status;type:enum('active', 'blocked');not null;default:'active'"`
	Role     Role           `gorm:"column:Role;type:enum('manager', 'customer');not null;default:'customer';index"`
	Phone    sql.NullString `gorm:"column:Phone;type:varchar(20)"`
	Avatar   sql.NullString `gorm:"column:Avatar;type:varchar(255)"`
}

func (u *User) TableName() string {
	return "User"
}

type CreateUserPayload struct {
	FullName string  `json:"fullName" validate:"required,max=255"`
	Email    string  `json:"email" validate:"required,email,max=255"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,max=20,phone_format"`
}

type UserResponse struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	RestaurantName string `json:"restaurantName,omitempty"`
	Avatar         string `json:"avatar,omitempty"`
}

func (payload *CreateUserPayload) ToUser(Role Role) *User {
	ID, _ := uuid.NewV7()
	return &User{
		BaseModel: BaseModel{
			ID: ID,
		},
		FullName: payload.FullName,
		Email:    payload.Email,
		Role:     Role,
		Phone:    sql.NullString{String: *payload.Phone, Valid: payload.Phone != nil},
	}
}

func (user *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
		Avatar:   user.Avatar.String,
	}
}
