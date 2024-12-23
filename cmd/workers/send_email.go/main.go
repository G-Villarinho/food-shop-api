package main

import (
	"context"
	"log"
	"time"

	"github.com/G-Villarinho/level-up-api/client"
	"github.com/G-Villarinho/level-up-api/config"
	"github.com/G-Villarinho/level-up-api/database"
	"github.com/G-Villarinho/level-up-api/internal"
	"github.com/G-Villarinho/level-up-api/models"
	"github.com/G-Villarinho/level-up-api/services"
	"github.com/G-Villarinho/level-up-api/services/email"
	"github.com/G-Villarinho/level-up-api/templates"
	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
)

func main() {
	config.ConfigureLogger()
	config.LoadEnvironments()

	di := internal.NewDi()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	redisClient, err := database.NewRedisConnection(ctx)
	if err != nil {
		log.Fatal("error to connect to redis: ", err)
	}

	internal.Provide(di, func(d *internal.Di) (*redis.Client, error) {
		return redisClient, nil
	})

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

	internal.Provide(di, client.NewMailtrapClient)
	internal.Provide(di, services.NewQueueService)
	internal.Provide(di, email.NewEmailService)
	internal.Provide(di, templates.NewTemplateService)

	emailService, err := internal.Invoke[email.EmailService](di)
	if err != nil {
		log.Fatal("error to create like service: ", err)
	}

	queueService, err := internal.Invoke[services.QueueService](di)
	if err != nil {
		log.Fatal("error to create queue service: ", err)
	}

	for {
		messages, err := queueService.Consume(services.QueueSendEmail)
		if err != nil {
			log.Fatal("error to consume message from queue: ", err)
		}

		for message := range messages {
			var task models.EmailQueueTask
			if err := jsoniter.Unmarshal(message, &task); err != nil {
				log.Println("error unmarshalling email task: ", err)
				continue
			}

			if err := emailService.SendEmail(context.Background(), task); err != nil {
				log.Println("error sending email: ", err)
				continue
			}

			log.Println("email sent successfully")
		}
	}

}
