package router

import (
	"log"

	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupUserRoutes(e *echo.Echo, di *internal.Di) {
	userHandler, err := internal.Invoke[handler.UserHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/users")

	group.POST("", userHandler.CreateUser)
	group.GET("/me", userHandler.GetUser, middleware.EnsureAuthenticated(di))
}
