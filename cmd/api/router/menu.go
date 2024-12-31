package router

import (
	"log"

	"github.com/G-Villarinho/food-shop-api/cmd/api/handler"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/middleware"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/labstack/echo/v4"
)

func setupMenuRoutes(e *echo.Echo, di *internal.Di) {
	menuHandler, err := internal.Invoke[handler.MenuHandler](di)
	if err != nil {
		log.Fatal("error to create menu handler: ", err)
	}

	group := e.Group("/v1/menus", middleware.EnsureAuthenticated(di))
	group.PUT("", menuHandler.UpdateMenu, middleware.EnsurePermission(models.UpdateMenuPermission))
}
