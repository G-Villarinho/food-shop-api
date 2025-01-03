package router

import (
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, di *internal.Di) {
	setupUserRoutes(e, di)
	setupAuthRoutes(e, di)
	setupRestaurantRoutes(e, di)
	setupOrderRoutes(e, di)
	setupProductRoutes(e, di)
	setupEvaluationRoutes(e, di)
	setupMenuRoutes(e, di)
	setupMetricsRouter(e, di)
}
