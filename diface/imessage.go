package diface

// IMessage ...
type IMessage interface {
	GetDataLen() uint32
	GetMsgID() uint32
	GetData() []byte

	SetDataLen(uint32)
	SetMsgID(uint32)
	SetData([]byte)
}
