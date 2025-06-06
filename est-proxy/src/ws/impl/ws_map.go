package impl

import (
	"est-proxy/src/ws"
	"github.com/google/uuid"
	"sync"
)

type ConnectionsMap struct {
	connections map[uuid.UUID][]ws.Connection
	lock        sync.RWMutex
}

func NewConnectionsMap() *ConnectionsMap {
	return &ConnectionsMap{
		connections: make(map[uuid.UUID][]ws.Connection),
		lock:        sync.RWMutex{},
	}
}

func (connMap *ConnectionsMap) Save(boardId uuid.UUID, conn ws.Connection) {
	connMap.lock.Lock()
	defer connMap.lock.Unlock()

	connections, ok := connMap.connections[boardId]
	if !ok {
		connections = make([]ws.Connection, 0)
	}
	connMap.connections[boardId] = append(connections, conn)
}

func (connMap *ConnectionsMap) Remove(boardId uuid.UUID, conn ws.Connection) {
	connMap.lock.Lock()
	defer connMap.lock.Unlock()

	connections, ok := connMap.connections[boardId]
	if ok {
		for i, c := range connections {
			if c == conn {
				connections = append(connections[:i], connections[i+1:]...)
				break
			}
		}
		connMap.connections[boardId] = connections
	}
}

func (connMap *ConnectionsMap) GetConnections(boardId uuid.UUID) []ws.Connection {
	connMap.lock.RLock()
	defer connMap.lock.RUnlock()

	connections, ok := connMap.connections[boardId]
	if !ok {
		return make([]ws.Connection, 0)
	}

	return connections
}
