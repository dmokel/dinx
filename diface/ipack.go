package diface

// IPack ...
type IPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack(data []byte) (IMessage, error)
}
