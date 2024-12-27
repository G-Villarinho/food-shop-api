package router

import (
	"log"

	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/middleware"
	"github.com/labstack/echo/v4"
)

func setupProductRoutes(e *echo.Echo, di *internal.Di) {
	productHandler, err := internal.Invoke[handler.ProductHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/products", middleware.EnsureAuthenticated(di))

	group.GET("/popular", productHandler.GetPopularProducts)
}
