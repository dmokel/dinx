package dnet

import "github.com/dmokel/dinx/diface"

// Request used to pack connection and data
type Request struct {
	connection diface.IConnection
	message    diface.IMessage
}

var _ diface.IRequest = &Request{}

// GetConnection ...
func (r *Request) GetConnection() diface.IConnection {
	return r.connection
}

// GetMessageID ...
func (r *Request) GetMessageID() uint32 {
	return r.message.GetMsgID()
}

// GetData ...
func (r *Request) GetData() []byte {
	return r.message.GetData()
}

// GetDataLength ...
func (r *Request) GetDataLength() uint32 {
	return r.message.GetDataLen()
}
