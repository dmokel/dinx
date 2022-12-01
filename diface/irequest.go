package diface

// IRequest ...
type IRequest interface {
	GetConnection() IConnection
	GetData() []byte
	GetDataLength() int
}
