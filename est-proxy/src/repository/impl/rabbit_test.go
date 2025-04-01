package impl

import (
	"testing"
	"time"

	"est-proxy/src/repository"
	"github.com/stretchr/testify/assert"
)

func TestRabbitRepositoryImpl_GetTopic(t *testing.T) {
	fakeCh := newFakeChannel()
	repo := &RabbitRepositoryImpl{
		declaredTopics: make(map[string]repository.RabbitTopic),
		channel:        fakeCh,
	}

	topicName := "test-topic"
	topic := repo.GetTopic(topicName)
	assert.NotNil(t, topic)
	_, exists := repo.declaredTopics[topicName]
	assert.True(t, exists, "Топик должен быть зарегистрирован в репозитории")

	msgReceived := make(chan []byte, 1)
	err := topic.Subscribe(func(msg []byte) {
		msgReceived <- msg
	})
	assert.NoError(t, err)

	testMsg := []byte("Hello, RabbitRepositoryImpl!")
	err = topic.Publish(testMsg)
	assert.NoError(t, err)
	select {
	case msg := <-msgReceived:
		assert.Equal(t, testMsg, msg)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Сообщение не получено в рамках RabbitRepositoryImpl.GetTopic")
	}
}

func TestRabbitRepositoryImpl_ReconnectTopics(t *testing.T) {
	repo := &RabbitRepositoryImpl{
		declaredTopics: make(map[string]repository.RabbitTopic),
	}
	fakeCh1 := newFakeChannel()
	repo.channel = fakeCh1

	topicName := "test-topic"
	topic := repo.GetTopic(topicName)

	msgReceived := make(chan []byte, 1)
	err := topic.Subscribe(func(msg []byte) {
		msgReceived <- msg
	})
	assert.NoError(t, err)

	testMsg1 := []byte("Сообщение до переподключения")
	err = topic.Publish(testMsg1)
	assert.NoError(t, err)
	select {
	case msg := <-msgReceived:
		assert.Equal(t, testMsg1, msg)
	case <-time.After(time.Second):
		t.Fatal("Сообщение до переподключения не получено")
	}

	fakeCh2 := newFakeChannel()
	repo.channel = fakeCh2
	repo.reconnectTopics()

	testMsg2 := []byte("Сообщение после переподключения")
	err = repo.GetTopic(topicName).Publish(testMsg2)
	assert.NoError(t, err)
	select {
	case msg := <-msgReceived:
		assert.Equal(t, testMsg2, msg)
	case <-time.After(time.Second):
		t.Fatal("Сообщение после переподключения не получено")
	}
}
