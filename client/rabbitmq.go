package client

import (
	"fmt"

	"github.com/G-Villarinho/food-shop-api/config"
	"github.com/G-Villarinho/food-shop-api/internal"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient interface {
	Connect() error
	Publish(queueName string, message []byte) error
	Consume(queueName string) (<-chan []byte, error)
	Disconnect() error
}

type rabbitMQClient struct {
	di         *internal.Di
	connection *amqp091.Connection
	channel    *amqp091.Channel
}

func NewRabbitMQClient(di *internal.Di) (RabbitMQClient, error) {
	return &rabbitMQClient{
		di: di,
	}, nil
}

func (r *rabbitMQClient) Connect() error {
	conn, err := amqp091.Dial(config.Env.RabbitMQURL)
	if err != nil {
		return err
	}
	r.connection = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	r.channel = channel

	return nil
}

func (r *rabbitMQClient) Disconnect() error {
	if r.channel != nil {
		err := r.channel.Close()
		if err != nil {
			return err
		}
	}

	if r.connection != nil {
		err := r.connection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *rabbitMQClient) Publish(queueName string, message []byte) error {
	if r.channel == nil {
		return fmt.Errorf("rabbitMQ channel is not initialized, ensure Connect() is called before Publish")
	}

	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = r.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)

	return err
}

func (r *rabbitMQClient) Consume(queueName string) (<-chan []byte, error) {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgs, err := r.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan []byte)
	go func() {
		for msg := range msgs {
			msgChan <- msg.Body
		}
		close(msgChan)
	}()

	return msgChan, nil
}
