package diface

// IConnectionManager ...
type IConnectionManager interface {
	Add(connection IConnection)
	Remove(connection IConnection)
	Num() int
	GetConnection(connID uint32) (IConnection, error)
	Clear()
}
