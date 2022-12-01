package dnet

import (
	"fmt"
	"net"

	"github.com/dmokel/dinx/diface"
)

type server struct {
	Name    string
	Network string
	IP      string
	Port    int

	Router diface.IRouter
}

var _ diface.IServer = &server{}

// NewServer used to create the server instance
func NewServer() diface.IServer {
	return &server{
		Name:    "default",
		Network: "tcp",
		IP:      "127.0.0.1",
		Port:    8999,

		Router: nil,
	}
}

func (s *server) Start() {
	go func() {
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

			conn := NewConnection(tcpConn, uint32(cid), s.Router)
			cid++
			go conn.Start()
		}
	}()
}

func (s *server) Stop() {}

func (s *server) Serve() {
	s.Start()

	fmt.Println("[Server] Dinx Server Serve...")
	select {}
}

// AddRouter ...
func (s *server) AddRouter(router diface.IRouter) {
	s.Router = router
	fmt.Println("[Server] add router")
}
