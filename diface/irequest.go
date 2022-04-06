package diface

// IRequest ...
type IRequest interface {
	GetConn() IConn
	GetData() []byte
	GetMsgID() uint32
}
