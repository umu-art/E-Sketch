package impl

import (
	"est-proxy/src/repository"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type TopicImpl struct {
	name    string
	channel *amqp091.Channel
}

func (topic *TopicImpl) Publish(message []byte) error {
	return topic.channel.Publish(
		topic.name,
		"*",
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

func (topic *TopicImpl) Subscribe(callback repository.Callback) error {
	queue, err := topic.channel.QueueDeclare(
		"est-proxy-"+uuid.New().String(),
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = topic.channel.QueueBind(
		queue.Name,
		"*",
		topic.name,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	consumer, err := topic.channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range consumer {
			callback(msg.Body)
			if err := msg.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %v", err)
			}
		}
	}()

	return nil
}
