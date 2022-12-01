package diface

import "net"

// IConnection ...
type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnectionID() uint32
	RemoteAddr() string
	Send(data []byte) error
}
