package dnet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/utils"
)

// Pack ...
type Pack struct{}

var _ diface.IPack = &Pack{}

// NewPack ...
func NewPack() diface.IPack {
	return &Pack{}
}

// GetHeadLen ...
func (p *Pack) GetHeadLen() uint32 {
	// dataLen uint32 -> 4 byte, msgID uint32 -> 4 byte
	return 8
}

// Pack ...
func (p *Pack) Pack(msg diface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Unpack ...
func (p *Pack) Unpack(headData []byte) (diface.IMessage, error) {
	reader := bytes.NewReader(headData)
	msg := &message{}

	if err := binary.Read(reader, binary.LittleEndian, &msg.dataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(reader, binary.LittleEndian, &msg.msgID); err != nil {
		return nil, err
	}

	if utils.GlobalIns.MaxPackageSize > 0 && msg.dataLen > utils.GlobalIns.MaxPackageSize {
		return nil, errors.New("data exceeds the maximum packet length")
	}

	return msg, nil
}
