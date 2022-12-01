package dnet

import (
	"fmt"
	"net"

	"github.com/dmokel/dinx/diface"
)

// Connection ...
type Connection struct {
	TCPConn      *net.TCPConn
	ConnectionID uint32
	isClosed     bool
	exitChan     chan bool

	Router diface.IRouter
}

// NewConnection ...
func NewConnection(conn *net.TCPConn, connectionID uint32, router diface.IRouter) *Connection {
	return &Connection{
		TCPConn:      conn,
		ConnectionID: connectionID,
		isClosed:     false,
		exitChan:     make(chan bool, 1),

		Router: router,
	}
}

func (c *Connection) startReader() {
	fmt.Printf("[Connection] connectionID = %d reader is runnning\n", c.ConnectionID)
	defer fmt.Printf("[Connection] connectionID = %d reader exit\n", c.ConnectionID)
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		cnt, err := c.TCPConn.Read(buf)
		if err != nil {
			fmt.Printf("[Connection] connectionID = %d failed to read bytes stream from tcp conn, err:%v\n", c.ConnectionID, err)
			continue
		}

		req := &Request{
			connection: c,
			data:       buf,
			cnt:        cnt,
		}

		go func(req diface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(req)
	}
}

// Start used to start the connection processing logic
func (c *Connection) Start() {
	fmt.Printf("[Connection] connectionID = %d is starting\n", c.ConnectionID)
	go c.startReader()
}

// Stop used to close a connection
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true
	c.TCPConn.Close()
	close(c.exitChan)
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

// Send used to get send byte stream data to client
func (c *Connection) Send(data []byte) error {
	if _, err := c.TCPConn.Write(data); err != nil {
		fmt.Printf("[Connection] connectionID = %d failed write data to client, err:%v\n", c.ConnectionID, err)
		return err
	}
	return nil
}
