package repository

import "github.com/rabbitmq/amqp091-go"

type RabbitRepository interface {
	GetTopic(name string) RabbitTopic
	Refresh()
	Close()
}

type (
	RabbitTopic interface {
		Publish(message []byte) error
		Subscribe(handler Callback) error
		Reconnect(channel RabbitChannel) error
	}
)

type RabbitChannel interface {
	IsClosed() bool
	Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error)
	QueueBind(name, key, exchange string, noWait bool, args amqp091.Table) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error
	Close() error
}

type Callback func(message []byte)
