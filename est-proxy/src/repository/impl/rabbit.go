package impl

import (
	"est-proxy/src/config"
	"est-proxy/src/repository"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/url"
)

type RabbitRepositoryImpl struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitRepositoryImpl() *RabbitRepositoryImpl {
	repositoryAddress := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RABBITMQ_USERNAME,
		url.QueryEscape(config.RABBITMQ_PASSWORD),
		config.RABBITMQ_HOST,
		config.RABBITMQ_PORT)

	conn, err := amqp091.Dial(repositoryAddress)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	return &RabbitRepositoryImpl{
		conn:    conn,
		channel: channel,
	}
}

func (r *RabbitRepositoryImpl) Close() {
	err := r.channel.Close()
	if err != nil {
		log.Printf("Failed to close channel: %v", err)
	}

	err = r.conn.Close()
	if err != nil {
		log.Printf("Failed to close connection: %v", err)
	}
}

func (r *RabbitRepositoryImpl) GetTopic(name string) repository.Topic {
	err := r.channel.ExchangeDeclare(
		name,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare topic exchange: %v", err)
	}

	return &TopicImpl{
		name:    name,
		channel: r.channel,
	}
}
