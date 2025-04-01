package impl

import (
	"est-proxy/src/ws"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChannelImpl_Listen(t *testing.T) {
	channel := NewChannelImpl()
	msgChan := make(chan []byte, 1)
	boardID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		channel.Listen(w, r, func(msg []byte, _ ws.Connection) {
			msgChan <- msg
		})
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?boardId=" + boardID.String()
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err, "WebSocket connection failed")
	defer conn.Close()

	time.Sleep(1 * time.Millisecond)
	require.Equal(t, 1, len(channel.GetConnectionsForBoard(boardID)))

	testMessage := []byte("integration test")
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, testMessage))

	select {
	case received := <-msgChan:
		assert.Equal(t, testMessage, received, "Should receive correct message")
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Message processing timeout")
	}
}
