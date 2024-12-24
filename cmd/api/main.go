package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/G-Villarinho/level-up-api/cache"
	"github.com/G-Villarinho/level-up-api/client"
	"github.com/G-Villarinho/level-up-api/cmd/api/handler"
	"github.com/G-Villarinho/level-up-api/cmd/api/router"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/database"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/repositories"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/G-Villarinho/level-up-api/services/email"
	"github.com/G-Villarinho/level-up-api/templates"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	e := echo.New()
	di := internal.NewDi()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${method} ${uri} ${status} ${latency_human} ${bytes_in} bytes_in ${bytes_out} bytes_out` + "\n",
	}))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{config.Env.FrontURL},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to database: ", err)
	}

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to redis: ", err)
	}

	rabbitMQClient, err := client.NewRabbitMQClient(di)
	if err != nil {
		log.Fatal("error initializing RabbitMQ client: ", err)
	}

	if err := rabbitMQClient.Connect(); err != nil {
		log.Fatal("error connecting to RabbitMQ: ", err)
	}

	defer func() {
		if err := rabbitMQClient.Disconnect(); err != nil {
			log.Println("error disconnecting from RabbitMQ:", err)
		}
	}()

	internal.Provide(di, func(d *internal.Di) (client.RabbitMQClient, error) {
		return rabbitMQClient, nil
	})

	internal.Provide(di, func(d *internal.Di) (*gorm.DB, error) {
		return db, nil
	})

	internal.Provide(di, func(d *internal.Di) (*redis.Client, error) {
		return redisClient, nil
	})

	internal.Provide(di, client.NewMailtrapClient)

	internal.Provide(di, handler.NewAuthHandler)
	internal.Provide(di, handler.NewOrderHandler)
	internal.Provide(di, handler.NewRestaurantHandler)
	internal.Provide(di, handler.NewUserHandler)

	internal.Provide(di, cache.NewRedisCache)
	internal.Provide(di, email.NewEmailService)
	internal.Provide(di, templates.NewTemplateService)

	internal.Provide(di, services.NewAuthService)
	internal.Provide(di, services.NewOrderItemService)
	internal.Provide(di, services.NewOrderService)
	internal.Provide(di, services.NewQueueService)
	internal.Provide(di, services.NewRestaurantService)
	internal.Provide(di, services.NewSessionService)
	internal.Provide(di, services.NewTokenService)
	internal.Provide(di, services.NewUserService)

	internal.Provide(di, repositories.NewOrderRepository)
	internal.Provide(di, repositories.NewProductRepository)
	internal.Provide(di, repositories.NewRestaurantRepository)
	internal.Provide(di, repositories.NewUserRepository)

	router.SetupRoutes(e, di)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Env.APIPort)))
}
