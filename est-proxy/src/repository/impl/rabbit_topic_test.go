package impl

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	esterrors "est-proxy/src/errors"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

type dummyAcknowledger struct{}

func (d *dummyAcknowledger) Ack(deliveryTag uint64, multiple bool) error {
	return nil
}

func (d *dummyAcknowledger) Nack(deliveryTag uint64, multiple, requeue bool) error {
	return nil
}

func (d *dummyAcknowledger) Reject(deliveryTag uint64, requeue bool) error {
	return nil
}

type fakeChannel struct {
	closed          bool
	mu              sync.Mutex
	consumers       map[string][]chan amqp091.Delivery
	queues          map[string]string
	errQueueDeclare error
	errQueueBind    error
	errConsume      error
}

func newFakeChannel() *fakeChannel {
	return &fakeChannel{
		closed:    false,
		consumers: make(map[string][]chan amqp091.Delivery),
		queues:    make(map[string]string),
	}
}

func (f *fakeChannel) IsClosed() bool {
	return f.closed
}

func (f *fakeChannel) Publish(exchange, key string, mandatory, immediate bool, msg amqp091.Publishing) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if chans, ok := f.consumers[exchange]; ok {
		for _, ch := range chans {
			delivery := amqp091.Delivery{
				Body:         msg.Body,
				Acknowledger: &dummyAcknowledger{},
			}
			go func(c chan amqp091.Delivery, d amqp091.Delivery) {
				c <- d
			}(ch, delivery)
		}
	}
	return nil
}

func (f *fakeChannel) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp091.Table) (amqp091.Queue, error) {
	if f.errQueueDeclare != nil {
		return amqp091.Queue{}, f.errQueueDeclare
	}
	q := amqp091.Queue{Name: name}
	f.mu.Lock()
	f.queues[name] = ""
	f.mu.Unlock()
	return q, nil
}

func (f *fakeChannel) QueueBind(name, key, exchange string, noWait bool, args amqp091.Table) error {
	if f.errQueueBind != nil {
		return f.errQueueBind
	}
	f.mu.Lock()
	f.queues[name] = exchange
	f.mu.Unlock()
	return nil
}

func (f *fakeChannel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp091.Table) (<-chan amqp091.Delivery, error) {
	if f.errConsume != nil {
		return nil, f.errConsume
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	exchange, ok := f.queues[queue]
	if !ok {
		return nil, amqp091.ErrClosed
	}
	ch := make(chan amqp091.Delivery, 10)
	f.consumers[exchange] = append(f.consumers[exchange], ch)
	return ch, nil
}

func (f *fakeChannel) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp091.Table) error {
	return nil
}

func (f *fakeChannel) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.closed {
		return errors.New("channel already closed")
	}

	f.closed = true

	for _, chans := range f.consumers {
		for _, ch := range chans {
			close(ch)
		}
	}

	f.consumers = make(map[string][]chan amqp091.Delivery)
	f.queues = make(map[string]string)

	return nil
}

func TestTopicImpl_Subscribe(t *testing.T) {
	t.Run("Нормальная подписка", func(t *testing.T) {
		fakeCh := newFakeChannel()
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh)

		msgReceived := make(chan []byte, 1)
		err := topic.Subscribe(func(msg []byte) {
			msgReceived <- msg
		})
		assert.NoError(t, err)

		expectedMsg := []byte("test message")
		publishing := amqp091.Publishing{
			Body: expectedMsg,
		}
		err = fakeCh.Publish(
			topicName,
			"",
			false,
			false,
			publishing,
		)
		assert.NoError(t, err)

		select {
		case msg := <-msgReceived:
			assert.Equal(t, expectedMsg, msg, "Полученное сообщение не совпадает с ожидаемым")
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout: сообщение не получено вовремя")
		}
	})

	t.Run("Ошибка QueueDeclare", func(t *testing.T) {
		fakeCh := newFakeChannel()
		fakeCh.errQueueDeclare = errors.New("ошибка queue declare")
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh)

		err := topic.Subscribe(func(msg []byte) {})
		assert.Error(t, err)
		assert.Equal(t, "ошибка queue declare", err.Error())
	})

	t.Run("Ошибка QueueBind", func(t *testing.T) {
		fakeCh := newFakeChannel()
		fakeCh.errQueueBind = errors.New("ошибка queue bind")
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh)

		err := topic.Subscribe(func(msg []byte) {})
		assert.Error(t, err)
		assert.Equal(t, "ошибка queue bind", err.Error())
	})

	t.Run("Ошибка Consume", func(t *testing.T) {
		fakeCh := newFakeChannel()
		fakeCh.errConsume = errors.New("ошибка consume")
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh)

		err := topic.Subscribe(func(msg []byte) {})
		assert.Error(t, err)
		assert.Equal(t, "ошибка consume", err.Error())
	})
}

func TestTopicImpl_Publish(t *testing.T) {
	t.Run("Нормальная отправка", func(t *testing.T) {
		fakeCh := newFakeChannel()
		topicName := "test"
		publisher := NewTopicImpl(topicName, fakeCh)

		queue, _ := fakeCh.QueueDeclare("queue1", false, false, true, false, nil)
		_ = fakeCh.QueueBind(queue.Name, "*", topicName, false, nil)
		consumer, _ := fakeCh.Consume(queue.Name, "", false, false, false, false, nil)

		testMsg := []byte("Привет, Publish!")
		err := publisher.Publish(testMsg)
		assert.NoError(t, err)

		select {
		case delivery := <-consumer:
			assert.Equal(t, testMsg, delivery.Body)
		case <-time.After(time.Millisecond):
			t.Fatal("Опубликованное сообщение не получено")
		}
	})

	t.Run("Отправка с закрытым каналом", func(t *testing.T) {
		fakeCh := newFakeChannel()
		fakeCh.closed = true
		topicName := "test"
		publisher := NewTopicImpl(topicName, fakeCh)

		err := publisher.Publish([]byte("Сообщение"))
		assert.Error(t, err)
		assert.Equal(t, esterrors.ErrRabbitChannelClosed.Error(), err.Error())
	})
}

func TestTopicImpl_Reconnect(t *testing.T) {
	t.Run("Успешное переподключение", func(t *testing.T) {
		fakeCh1 := newFakeChannel()
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh1)

		msgReceived := make(chan []byte, 1)
		err := topic.Subscribe(func(msg []byte) {
			msgReceived <- msg
		})
		assert.NoError(t, err)

		testMsg1 := []byte("Сообщение до переподключения")
		err = topic.Publish(testMsg1)
		assert.NoError(t, err)
		select {
		case m := <-msgReceived:
			assert.Equal(t, testMsg1, m)
		case <-time.After(time.Second):
			t.Fatal("Сообщение до переподключения не получено")
		}

		fakeCh2 := newFakeChannel()
		err = topic.Reconnect(fakeCh2)
		assert.NoError(t, err)

		testMsg2 := []byte("Сообщение после переподключения")
		err = topic.Publish(testMsg2)
		assert.NoError(t, err)
		select {
		case m := <-msgReceived:
			assert.Equal(t, testMsg2, m)
		case <-time.After(time.Second):
			t.Fatal("Сообщение после переподключения не получено")
		}
	})

	t.Run("Ошибка переподключения при сбое подписки", func(t *testing.T) {
		fakeCh1 := newFakeChannel()
		topicName := "test"
		topic := NewTopicImpl(topicName, fakeCh1)

		err := topic.Subscribe(func(msg []byte) {})
		assert.NoError(t, err)

		fakeCh2 := newFakeChannel()
		fakeCh2.errQueueDeclare = errors.New("ошибка queue declare при переподключении")
		err = topic.Reconnect(fakeCh2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ошибка queue declare при переподключении")
	})
}

func TestTopicImpl_ManyToMany(t *testing.T) {
	t.Run("Многопоточность подписчиков и издателей", func(t *testing.T) {
		fakeCh := newFakeChannel()
		topicName := "test-mtm"

		numSubscribers := 5
		numPublishers := 5
		messagesPerPublisher := 10

		var subscribers []chan []byte
		for i := 0; i < numSubscribers; i++ {
			topic := NewTopicImpl(topicName, fakeCh)
			ch := make(chan []byte, numPublishers*messagesPerPublisher)
			err := topic.Subscribe(func(msg []byte) {
				ch <- msg
			})
			assert.NoError(t, err)
			subscribers = append(subscribers, ch)
		}

		var publishers []*TopicImpl
		for i := 0; i < numPublishers; i++ {
			publishers = append(publishers, NewTopicImpl(topicName, fakeCh))
		}

		var wg sync.WaitGroup
		for i, pub := range publishers {
			wg.Add(1)
			go func(pub *TopicImpl, publisherIndex int) {
				defer wg.Done()
				for j := 0; j < messagesPerPublisher; j++ {
					msg := []byte(fmt.Sprintf("Publisher %d message %d", publisherIndex, j))
					err := pub.Publish(msg)
					assert.NoError(t, err)
				}
			}(pub, i)
		}
		wg.Wait()

		expectedTotal := numPublishers * messagesPerPublisher
		for i, sub := range subscribers {
			var received []string
			timeout := time.After(2 * time.Second)
			for len(received) < expectedTotal {
				select {
				case msg := <-sub:
					received = append(received, string(msg))
				case <-timeout:
					t.Fatalf("Подписчик %d получил %d сообщений, ожидается %d", i, len(received), expectedTotal)
				}
			}
			assert.Equal(t, expectedTotal, len(received))
		}
	})
}
