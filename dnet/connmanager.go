package dnet

import (
	"errors"
	"sync"

	"github.com/dmokel/dinx/diface"
)

// ConnectionManager ...
type connectionManager struct {
	connections map[uint32]diface.IConnection
	lock        sync.RWMutex
}

// NewConnectionManager ...
func newConnectionManager() diface.IConnectionManager {
	return &connectionManager{
		connections: make(map[uint32]diface.IConnection),
	}
}

var _ diface.IConnectionManager = &connectionManager{}

// Add ...
func (cm *connectionManager) Add(connection diface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.connections[connection.GetConnectionID()] = connection
}

// Remove ...
func (cm *connectionManager) Remove(connection diface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.connections, connection.GetConnectionID())
}

// Num ...
func (cm *connectionManager) Num() int {
	return len(cm.connections)
}

// GetConnection ...
func (cm *connectionManager) GetConnection(connID uint32) (diface.IConnection, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if connection, ok := cm.connections[connID]; ok {
		return connection, nil
	}
	return nil, errors.New("not match any connection")
}

// Clear ...
func (cm *connectionManager) Clear() {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for connID, connection := range cm.connections {
		connection.Close()
		delete(cm.connections, connID)
	}
}
