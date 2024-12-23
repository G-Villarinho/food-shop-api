package services

import (
	"fmt"

	"github.com/G-Villarinho/level-up-api/client"
	"github.com/G-Villarinho/level-up-api/internal"
)

const (
	QueueSendEmail = "send_email_queue"
)

type QueueService interface {
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan []byte, error)
}

type queueService struct {
	di             *internal.Di
	rabbitMQClient client.RabbitMQClient
}

func NewQueueService(di *internal.Di) (QueueService, error) {
	rabbitMQClient, err := internal.Invoke[client.RabbitMQClient](di)
	if err != nil {
		return nil, err
	}

	return &queueService{
		di:             di,
		rabbitMQClient: rabbitMQClient,
	}, nil
}

func (q *queueService) Publish(queueName string, message []byte) error {
	if err := q.rabbitMQClient.Publish(queueName, message); err != nil {
		return fmt.Errorf("publishing message to queue: %w", err)
	}

	return nil
}

func (q *queueService) Consume(queueName string) (<-chan []byte, error) {
	messages, err := q.rabbitMQClient.Consume(queueName)
	if err != nil {
		return nil, fmt.Errorf("consuming message from queue: %w", err)
	}

	return messages, nil
}
