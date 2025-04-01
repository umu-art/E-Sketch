package ws

type Connection interface {
	ReadMessage() ([]byte, error)
	WriteMessage(message []byte) error
	Close() error
}
