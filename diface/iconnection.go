package diface

import "net"

// IConnection ...
type IConnection interface {
	Start()
	Stop()
	GetTCPConn() *net.TCPConn
	GetConnectionID() uint32
	RemoteAddr() string
	SendMsg(msgID uint32, data []byte) error
}
