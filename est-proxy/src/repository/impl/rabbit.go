package impl

import (
	"est-proxy/src/config"
	"est-proxy/src/repository"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/url"
	"time"
)

type RabbitRepositoryImpl struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitRepositoryImpl() *RabbitRepositoryImpl {
	var rabbitRepo RabbitRepositoryImpl
	if err := rabbitRepo.connect(); err != nil {
		log.Fatalf("%v", err)
	}
	return &rabbitRepo
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

func (r *RabbitRepositoryImpl) Refresh() {
	for {
		<-r.conn.NotifyClose(make(chan *amqp091.Error))
		failedAttempts := 0
		for {
			time.Sleep(5 * time.Second)
			if err := r.connect(); err == nil {
				break
			}

			if failedAttempts++; failedAttempts > 5 {
				log.Fatalf("Failed to reconnect to RabbitMQ")
			}
		}
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

func (r *RabbitRepositoryImpl) connect() error {
	repositoryAddress := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RABBITMQ_USERNAME,
		url.QueryEscape(config.RABBITMQ_PASSWORD),
		config.RABBITMQ_HOST,
		config.RABBITMQ_PORT)

	conn, err := amqp091.Dial(repositoryAddress)
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open channel: %v", err)
	}

	r.conn = conn
	r.channel = channel
	return nil
}
