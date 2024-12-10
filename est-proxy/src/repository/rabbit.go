package repository

type RabbitRepository interface {
	GetTopic(name string) Topic
	Close()
}

type Topic interface {
	Publish(message []byte) error
	Subscribe(handler Callback) error
}

type Callback func(message []byte)
