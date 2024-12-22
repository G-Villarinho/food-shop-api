package models

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound       = errors.New("user not found in the database")
	ErrEmailAlreadyExists = errors.New("email already exists in the database")
)

type Status string

const (
	Active  Status = "Active"
	Blocked Status = "Blocked"
)

type User struct {
	BaseModel
	FullName string         `gorm:"column:FullName;type:varchar(255);not null"`
	Email    string         `gorm:"column:Email;type:varchar(255);not null;unique"`
	Status   Status         `gorm:"column:Status;type:enum('Active', 'Blocked');not null;default:'Active'"`
	Avatar   sql.NullString `gorm:"column:Avatar;type:varchar(255)"`
	XP       int            `gorm:"column:XP;type:int;not null;default:0"`
}

func (u *User) TableName() string {
	return "User"
}

type CreateUserPayload struct {
	FullName string `json:"full_name" validate:"required,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
}

type UserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar,omitempty"`
	XP       int    `json:"xp"`
}

func (payload *CreateUserPayload) ToUser() *User {
	ID, _ := uuid.NewV7()
	return &User{
		BaseModel: BaseModel{
			ID: ID,
		},
		FullName: payload.FullName,
		Email:    payload.Email,
		XP:       0,
	}
}

func (user *User) ToUserResponse() *UserResponse {
	return &UserResponse{
		ID:       user.ID.String(),
		FullName: user.FullName,
		Email:    user.Email,
		Avatar:   user.Avatar.String,
		XP:       user.XP,
	}
}
