package dnet

import "github.com/dmokel/dinx/diface"

// Request used to pack connection and data
type Request struct {
	connection diface.IConnection
	data       []byte
	cnt        int
}

// GetConnection ...
func (r *Request) GetConnection() diface.IConnection {
	return r.connection
}

// GetData ...
func (r *Request) GetData() []byte {
	return r.data
}

// GetDataLength ...
func (r *Request) GetDataLength() int {
	return r.cnt
}
