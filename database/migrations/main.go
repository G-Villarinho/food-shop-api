package main

import (
	"context"
	"log"
	"time"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/database"
	"github.com/G-Villarinho/food-shop-api/models"
)

func main() {
	config.LoadEnvironments()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := database.NewMysqlConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to mysql: ", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Restaurant{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Evaluation{},
	); err != nil {
		log.Fatal("error to migrate: ", err)
	}

	log.Println("Migration executed successfully")

}
