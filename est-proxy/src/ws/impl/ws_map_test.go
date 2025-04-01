package impl

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createTestWSConn(t *testing.T) *websocket.Conn {
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
	}))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)

	return conn
}

func TestConnectionsMap_Save(t *testing.T) {
	connMap := NewConnectionsMap()
	boardID := uuid.New()

	conn1 := NewConnectionImpl(createTestWSConn(t))
	defer conn1.Close()
	conn2 := NewConnectionImpl(createTestWSConn(t))
	defer conn2.Close()

	t.Run("Первый подключение к доске", func(t *testing.T) {
		connMap.Save(boardID, conn1)
		result := connMap.GetConnections(boardID)
		assert.Len(t, result, 1)
		assert.Contains(t, result, conn1)
	})

	t.Run("Второй подключение к той же доске", func(t *testing.T) {
		connMap.Save(boardID, conn2)
		result := connMap.GetConnections(boardID)
		assert.Len(t, result, 2)
		assert.Contains(t, result, conn2)
	})
}

func TestConnectionsMap_Remove(t *testing.T) {
	connMap := NewConnectionsMap()
	boardID := uuid.New()
	conn := NewConnectionImpl(createTestWSConn(t))
	defer conn.Close()

	connMap.Save(boardID, conn)

	t.Run("Существует подключение", func(t *testing.T) {
		connMap.Remove(boardID, conn)
		result := connMap.GetConnections(boardID)
		assert.Empty(t, result)
	})

	t.Run("Отсутствует подключение", func(t *testing.T) {
		connMap.Remove(boardID, conn) // No panic
	})
}

func TestConnectionsMap_GetConnections(t *testing.T) {
	connMap := NewConnectionsMap()
	boardID := uuid.New()
	otherBoardID := uuid.New()

	conn := NewConnectionImpl(createTestWSConn(t))
	defer conn.Close()

	t.Run("Существующая доска", func(t *testing.T) {
		connMap.Save(boardID, conn)
		result := connMap.GetConnections(boardID)
		assert.Len(t, result, 1)
		assert.Equal(t, conn, result[0])
	})

	t.Run("Необработанная доска", func(t *testing.T) {
		result := connMap.GetConnections(otherBoardID)
		assert.Empty(t, result)
	})
}
