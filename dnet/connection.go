package dnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/utils"
)

// Connection ...
type Connection struct {
	Server       diface.IServer
	TCPConn      *net.TCPConn
	ConnectionID uint32
	isClosed     bool
	msgChan      chan []byte
	exitChan     chan bool

	propertys    map[string]interface{}
	propertyLock sync.RWMutex

	RouterGroup diface.IRouterGroup
}

var _ diface.IConnection = &Connection{}

// NewConnection ...
func NewConnection(server diface.IServer, conn *net.TCPConn, connectionID uint32, routerGroup diface.IRouterGroup) *Connection {
	c := &Connection{
		Server:       server,
		TCPConn:      conn,
		ConnectionID: connectionID,
		isClosed:     false,
		msgChan:      make(chan []byte),
		exitChan:     make(chan bool, 1),

		propertys: make(map[string]interface{}),

		RouterGroup: routerGroup,
	}

	c.Server.GetConnectionManager().Add(c)
	return c
}

func (c *Connection) startReader() {
	fmt.Printf("[Connection] connectionID = %d reader is running\n", c.ConnectionID)
	defer fmt.Printf("[Connection] connectionID = %d reader exit\n", c.ConnectionID)
	defer c.Close()

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
	c.Server.CallOnConnStart(c)
	go c.startReader()
	go c.startWriter()
}

// Close used to close a connection
func (c *Connection) Close() {
	if c.isClosed == true {
		return
	}

	c.Server.CallOnConnClose(c)
	c.isClosed = true
	c.TCPConn.Close()
	c.Server.GetConnectionManager().Remove(c)
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

// SetProperty ...
func (c *Connection) SetProperty(key string, value interface{}) error {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.propertys[key]; ok {
		return errors.New("duplicate key")
	}
	c.propertys[key] = value
	return nil
}

// GetProperty ...
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if v, ok := c.propertys[key]; ok {
		return v, nil
	}
	return nil, errors.New("not match any key")
}

// RemoveProperty ...
func (c *Connection) RemoveProperty(key string) error {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if _, ok := c.propertys[key]; !ok {
		return errors.New("not match any key")
	}
	delete(c.propertys, key)
	return nil
}
