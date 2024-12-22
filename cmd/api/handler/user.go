package handler

import (
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
}

type userHandler struct {
	di          *internal.Di
	userService services.UserService
}

func NewUserHandler(di *internal.Di) (UserHandler, error) {
	userService, err := internal.Invoke[services.UserService](di)
	if err != nil {
		return nil, err
	}

	return &userHandler{
		di:          di,
		userService: userService,
	}, nil
}

func (u *userHandler) CreateUser(c echo.Context) error {
	panic("unimplemented")
}
