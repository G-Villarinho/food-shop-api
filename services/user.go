package services

import (
	"context"

	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, payload models.CreateUserPayload) error
}

type userService struct {
	di             *internal.Di
	userRepository repositories.UserRepository
}

func NewUserService(di *internal.Di) (UserService, error) {
	userRepository, err := internal.Invoke[repositories.UserRepository](di)
	if err != nil {
		return nil, err
	}

	return &userService{
		di:             di,
		userRepository: userRepository,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, payload models.CreateUserPayload) error {
	panic("unimplemented")
}
