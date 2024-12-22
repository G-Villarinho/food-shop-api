package main

import (
	"context"
	"log"
	"time"

	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/database"
	"github.com/G-Villarinho/level-up-api/models"
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
	); err != nil {
		log.Fatal("error to migrate: ", err)
	}

	log.Println("Migration executed successfully")

}
