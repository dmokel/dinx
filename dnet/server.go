package dnet

import (
	"fmt"
	"net"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/utils"
)

type server struct {
	Name    string
	Network string
	IP      string
	Port    int
	Version string

	RouterGroup diface.IRouterGroup
	connManager diface.IConnectionManager

	onConnstart diface.ConnHookFunc
	onConnClose diface.ConnHookFunc
}

var _ diface.IServer = &server{}

// NewServer used to create the server instance
func NewServer() diface.IServer {
	return &server{
		Name:    utils.GlobalIns.Name,
		Network: utils.GlobalIns.Network,
		IP:      utils.GlobalIns.IP,
		Port:    utils.GlobalIns.Port,
		Version: utils.GlobalIns.Version,

		RouterGroup: NewRouterGroup(),
		connManager: NewConnectionManager(),
	}
}

func (s *server) Start() {
	go func() {
		s.RouterGroup.StartWorkerPool()

		localAddr, err := net.ResolveTCPAddr(s.Network, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Server] failed to resolve tcp addr, err: ", err)
			return
		}

		listener, err := net.ListenTCP(s.Network, localAddr)
		if err != nil {
			fmt.Println("[Server] failed to create tcp listener, err: ", err)
			return
		}
		fmt.Printf("[Server] success listening tcp connection at %s:%d\n", s.IP, s.Port)

		cid := 0

		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[Server] failed to accept tcp connection, err: ", err)
				continue
			}

			if s.connManager.Num() >= utils.GlobalIns.MaxConn {
				fmt.Println("to many connections")
				tcpConn.Close()
				continue
			}

			conn := NewConnection(s, tcpConn, uint32(cid), s.RouterGroup)
			cid++
			go conn.Start()
		}
	}()
}

func (s *server) Stop() {
	s.connManager.Clear()
}

func (s *server) Serve() {
	s.Start()

	fmt.Printf("[Server] %s, Version %s is serving...\n", s.Name, s.Version)
	select {}
}

func (s *server) GetConnectionManager() diface.IConnectionManager {
	return s.connManager
}

// AddRouter ...
func (s *server) AddRouter(msgID uint32, router diface.IRouter) {
	s.RouterGroup.AddRouter(msgID, router)
}

func (s *server) SetOnConnStart(onConnStart diface.ConnHookFunc) {
	s.onConnstart = onConnStart
}

func (s *server) SetOnConnClose(onConnClose diface.ConnHookFunc) {
	s.onConnClose = onConnClose
}

func (s *server) CallOnConnStart(connection diface.IConnection) {
	if s.onConnstart != nil {
		s.onConnstart(connection)
	}
}

func (s *server) CallOnConnClose(connection diface.IConnection) {
	if s.onConnClose != nil {
		s.onConnClose(connection)
	}
}
