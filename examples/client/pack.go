package main

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/dmokel/dinx/utils"
)

// Pack ...
type pack struct{}

// NewPack ...
func newPack() *pack {
	return &pack{}
}

// GetHeadLen ...
func (p *pack) GetHeadLen() uint32 {
	// dataLen uint32 -> 4 byte, msgID uint32 -> 4 byte
	return 8
}

// Pack ...
func (p *pack) Pack(msg *message) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})

	if err := binary.Write(buffer, binary.LittleEndian, msg.dataLen); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.msgID); err != nil {
		return nil, err
	}

	if err := binary.Write(buffer, binary.LittleEndian, msg.data); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

type message struct {
	dataLen uint32
	msgID   uint32
	data    []byte
}

// Unpack ...
func (p *pack) Unpack(headData []byte) (*message, error) {
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
