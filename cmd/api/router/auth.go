package router

import (
	"log"

	"github.com/G-Villarinho/food-shop-api/cmd/api/handler"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/labstack/echo/v4"
)

func setupAuthRoutes(e *echo.Echo, di *internal.Di) {
	authHandler, err := internal.Invoke[handler.AuthHandler](di)
	if err != nil {
		log.Fatal("error to create auth handler: ", err)
	}

	group := e.Group("/v1/auth")

	group.POST("/sign-in", authHandler.SignIn)
	group.GET("/link", authHandler.VeryfyMagicLink)
}
