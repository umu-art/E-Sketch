package repository

import "github.com/rabbitmq/amqp091-go"

type RabbitRepository interface {
	GetTopic(name string) Topic
	Refresh()
	Close()
}

type Topic interface {
	Publish(message []byte) error
	Subscribe(handler Callback) error
	Reconnect(channel *amqp091.Channel) error
}

type Callback func(message []byte)
