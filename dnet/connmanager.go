package dnet

import (
	"errors"
	"sync"

	"github.com/dmokel/dinx/diface"
)

// ConnectionManager ...
type ConnectionManager struct {
	connections map[uint32]diface.IConnection
	lock        sync.RWMutex
}

// NewConnectionManager ...
func NewConnectionManager() diface.IConnectionManager {
	return &ConnectionManager{
		connections: make(map[uint32]diface.IConnection),
	}
}

var _ diface.IConnectionManager = &ConnectionManager{}

// Add ...
func (cm *ConnectionManager) Add(connection diface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.connections[connection.GetConnectionID()] = connection
}

// Remove ...
func (cm *ConnectionManager) Remove(connection diface.IConnection) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.connections, connection.GetConnectionID())
}

// Num ...
func (cm *ConnectionManager) Num() int {
	return len(cm.connections)
}

// GetConnection ...
func (cm *ConnectionManager) GetConnection(connID uint32) (diface.IConnection, error) {
	cm.lock.RLock()
	defer cm.lock.RUnlock()

	if connection, ok := cm.connections[connID]; ok {
		return connection, nil
	}
	return nil, errors.New("not match any connection")
}

// Clear ...
func (cm *ConnectionManager) Clear() {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for connID, connection := range cm.connections {
		connection.Close()
		delete(cm.connections, connID)
	}
}
