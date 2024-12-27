package router

import (
	"log"

	"github.com/G-Villarinho/food-shop-api/cmd/api/handler"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/G-Villarinho/food-shop-api/middleware"
	"github.com/G-Villarinho/food-shop-api/models"
	"github.com/labstack/echo/v4"
)

func setupEvaluationRoutes(e *echo.Echo, di *internal.Di) {
	evaluationHandler, err := internal.Invoke[handler.EvaluationHandler](di)
	if err != nil {
		log.Fatal("error to create evaluation handler: ", err)
	}

	group := e.Group("/v1/evaluations", middleware.EnsureAuthenticated(di))

	group.POST("", evaluationHandler.CreateEvaluation, middleware.EnsurePermission(models.CreateEvaluationPermission))
	group.GET("", evaluationHandler.GetEvaluations, middleware.EnsurePermission(models.ListEvaluationsPermission))
}
