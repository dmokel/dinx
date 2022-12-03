package dnet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/utils"
)

// Connection ...
type Connection struct {
	TCPConn      *net.TCPConn
	ConnectionID uint32
	isClosed     bool
	msgChan      chan []byte
	exitChan     chan bool

	RouterGroup diface.IRouterGroup
}

var _ diface.IConnection = &Connection{}

// NewConnection ...
func NewConnection(conn *net.TCPConn, connectionID uint32, routerGroup diface.IRouterGroup) *Connection {
	return &Connection{
		TCPConn:      conn,
		ConnectionID: connectionID,
		isClosed:     false,
		msgChan:      make(chan []byte),
		exitChan:     make(chan bool, 1),

		RouterGroup: routerGroup,
	}
}

func (c *Connection) startReader() {
	fmt.Printf("[Connection] connectionID = %d reader is running\n", c.ConnectionID)
	defer fmt.Printf("[Connection] connectionID = %d reader exit\n", c.ConnectionID)
	defer c.Stop()

	pack := NewPack()
	for {
		headBuf := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(c.TCPConn, headBuf); err != nil {
			fmt.Printf("[Connection] connectionID = %d failed to read message head, err:%v\n", c.ConnectionID, err)
			break
		}

		msg, err := pack.Unpack(headBuf)
		if err != nil {
			fmt.Printf("[Connection] connectionID = %d failed to unpack message head, err:%v", c.ConnectionID, err)
			break
		}

		var dataBuf []byte
		if msg.GetDataLen() > 0 {
			dataBuf = make([]byte, msg.GetDataLen())
			if _, err = io.ReadFull(c.TCPConn, dataBuf); err != nil {
				fmt.Printf("[Connection] connectionID = %d failed to read message data, err:%v", c.ConnectionID, err)
				break
			}
		}
		msg.SetData(dataBuf)

		req := &Request{
			connection: c,
			message:    msg,
		}

		if utils.GlobalIns.WorkerPoolSize > 0 {
			c.RouterGroup.SendMsgToTaskQueue(req)
		} else {
			go c.RouterGroup.DoMessageRouter(req)
		}
	}
}

func (c *Connection) startWriter() {
	fmt.Printf("[Connection] connectionID = %d writer is running\n", c.ConnectionID)
	defer fmt.Printf("[Connection] connectionID = %d writer exit\n", c.ConnectionID)

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.TCPConn.Write(data); err != nil {
				fmt.Printf("[Connection] connectionID = %d failed to write data back to client, err:%v\n", c.ConnectionID, err)
				return
			}
		case <-c.exitChan:
			return
		}
	}
}

// Start used to start the connection processing logic
func (c *Connection) Start() {
	go c.startReader()
	go c.startWriter()
}

// Stop used to close a connection
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true
	c.TCPConn.Close()
	close(c.exitChan)
	close(c.msgChan)
}

// GetTCPConn used to get the low level tcp conn
func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.TCPConn
}

// GetConnectionID used to get the connection id
func (c *Connection) GetConnectionID() uint32 {
	return c.ConnectionID
}

// RemoteAddr used to get the connection's remote addr
func (c *Connection) RemoteAddr() string {
	return c.TCPConn.RemoteAddr().String()
}

// SendMsg used to get send byte data to client
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New(`connection closed`)
	}

	pack := NewPack()
	msg := &message{
		msgID:   msgID,
		dataLen: uint32(len(data)),
		data:    data,
	}
	buf, err := pack.Pack(msg)
	if err != nil {
		fmt.Printf("[Connection] connectionID = %d failed to pack msg to buffer, err:%v\n", c.ConnectionID, err)
		return err
	}

	c.msgChan <- buf
	return nil
}
