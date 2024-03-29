package diface

import "net"

// IConnection ...
type IConnection interface {
	Start()
	Close()
	GetTCPConn() *net.TCPConn
	GetConnectionID() uint32
	RemoteAddr() string
	SendMsg(msgID uint32, data []byte) error

	SetProperty(key string, value interface{}) error
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string) error
}
