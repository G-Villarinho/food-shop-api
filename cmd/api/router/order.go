package router

import (
	"log"

	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupOrderRoutes(e *echo.Echo, di *internal.Di) {
	orderHandler, err := internal.Invoke[handler.OrderHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/orders", middleware.EnsureAuthenticated(di))

	group.POST("", orderHandler.CreateOrder)
}