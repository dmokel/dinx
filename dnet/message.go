package dnet

import "github.com/dmokel/dinx/diface"

// message ...
type message struct {
	dataLen uint32
	msgID   uint32
	data    []byte
}

var _ diface.IMessage = &message{}

// NewMessage ...
func NewMessage() diface.IMessage {
	return &message{}
}

// GetDataLen ...
func (m *message) GetDataLen() uint32 {
	return m.dataLen
}

// GettMsgID ...
func (m *message) GetMsgID() uint32 {
	return m.msgID
}

// GetData ...
func (m *message) GetData() []byte {
	return m.data
}

// SetDataLen ...
func (m *message) SetDataLen(dataLen uint32) {
	m.dataLen = dataLen
}

// SetMsgID ...
func (m *message) SetMsgID(msgID uint32) {
	m.msgID = msgID
}

// SetData ...
func (m *message) SetData(data []byte) {
	m.data = data
}
