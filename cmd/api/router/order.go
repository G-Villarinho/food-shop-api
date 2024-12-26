package router

import (
	"log"

	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/middleware"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/labstack/echo/v4"
)

func setupOrderRoutes(e *echo.Echo, di *internal.Di) {
	orderHandler, err := internal.Invoke[handler.OrderHandler](di)
	if err != nil {
		log.Fatal("error to create user handler: ", err)
	}

	group := e.Group("/v1/orders", middleware.EnsureAuthenticated(di))

	group.GET("", orderHandler.GetOrders, middleware.EnsurePermission(models.ListOrdersPermission))
	group.PATCH("/:orderId/cancel", orderHandler.CancelOrder, middleware.EnsurePermission(models.CancelOrderPermission))
}
