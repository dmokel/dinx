package diface

// IRequest ...
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetMessageID() uint32
	GetDataLength() uint32
}
