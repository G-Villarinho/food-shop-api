package router

import (
	"log"

	"github.com/G-Villarinho/food-shop-api/cmd/api/handler"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/middleware"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/labstack/echo/v4"
)

func setupMetricsRouter(e *echo.Echo, di *internal.Di) {
	metricsRouter, err := internal.Invoke[handler.MetricsHandler](di)
	if err != nil {
		log.Fatal("error to create metrics handler: ", err)
	}

	group := e.Group("/v1/metrics", middleware.EnsureAuthenticated(di))

	group.GET("/orders/monthly-amount", metricsRouter.GetMonthlyMetrics, middleware.EnsurePermission(models.GetMonthlyMetricsPermission))
}
