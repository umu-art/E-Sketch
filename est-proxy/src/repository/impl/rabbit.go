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
	conn      *amqp091.Connection
	channel   *amqp091.Channel
	topicList []repository.Topic
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
		time.Sleep(5 * time.Second)
		if !r.conn.IsClosed() {
			continue
		}

		log.Println("RabbitMQ connection is closed, attempting to reconnect")
		failedAttempts := 0
		for {
			time.Sleep(5 * time.Second)
			log.Printf("Trying to reconnect to RabbitMQ")
			if err := r.connect(); err == nil {
				r.reconnectTopics()
				log.Println("Successfully reconnected to RabbitMQ")
				break
			}
			failedAttempts++
			if failedAttempts > 5 {
				log.Fatalf("Failed to reconnect to RabbitMQ after %d attempts", failedAttempts)
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

	topic := NewTopicImpl(name, r.channel)
	r.topicList = append(r.topicList, topic)

	return topic
}

func (r *RabbitRepositoryImpl) connect() error {
	repositoryAddress := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RABBITMQ_USERNAME,
		url.QueryEscape(config.RABBITMQ_PASSWORD),
		config.RABBITMQ_HOST,
		config.RABBITMQ_PORT)

	conn, err := amqp091.Dial(repositoryAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	if r.channel != nil && !r.channel.IsClosed() {
		err = r.channel.Close()
		if err != nil {
			log.Printf("Failed to close channel: %v", err)
		}
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %v", err)
	}

	r.conn = conn
	r.channel = channel
	return nil
}

func (r *RabbitRepositoryImpl) reconnectTopics() {
	for _, topic := range r.topicList {
		err := topic.Reconnect(r.channel)
		if err != nil {
			log.Fatalf("Failed to reconnect to topic %v", err)
		}
	}
}
