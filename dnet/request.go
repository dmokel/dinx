package dnet

import "github.com/dmokel/dinx/diface"

// Request used to pack connection and data
type request struct {
	connection diface.IConnection
	message    diface.IMessage
}

var _ diface.IRequest = &request{}

// GetConnection ...
func (r *request) GetConnection() diface.IConnection {
	return r.connection
}

// GetMessageID ...
func (r *request) GetMessageID() uint32 {
	return r.message.GetMsgID()
}

// GetData ...
func (r *request) GetData() []byte {
	return r.message.GetData()
}

// GetDataLength ...
func (r *request) GetDataLength() uint32 {
	return r.message.GetDataLen()
}
