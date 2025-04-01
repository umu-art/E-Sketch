package impl

import (
	esterrors "est-proxy/src/errors"
	"est-proxy/src/repository"
	"fmt"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type TopicImpl struct {
	name      string
	channel   repository.RabbitChannel
	callbacks []repository.Callback
}

func NewTopicImpl(name string, channel repository.RabbitChannel) *TopicImpl {
	return &TopicImpl{
		name:    name,
		channel: channel,
	}
}

func (topic *TopicImpl) Publish(message []byte) error {
	if topic.channel == nil || topic.channel.IsClosed() {
		return esterrors.ErrRabbitChannelClosed
	}
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

	topic.callbacks = append(topic.callbacks, callback)

	return nil
}

func (topic *TopicImpl) Reconnect(channel repository.RabbitChannel) error {
	topic.channel = channel
	for _, callback := range topic.callbacks {
		err := topic.Subscribe(callback)
		if err != nil {
			return fmt.Errorf("error subscribing to topic %s: %v", topic.name, err)
		}
	}

	return nil
}
