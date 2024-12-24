package router

import (
	"log"

	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/middleware"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/labstack/echo/v4"
)

func setupRestaurantRoutes(e *echo.Echo, di *internal.Di) {
	restaurantHandler, err := internal.Invoke[handler.RestaurantHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/restaurants")

	group.POST("", restaurantHandler.CreateRestaurant)
	group.POST("/:restaurantID/order", restaurantHandler.CreateOrder, middleware.EnsureAuthenticated(di), middleware.EnsurePermission(models.CreateOrderPermission))
}
