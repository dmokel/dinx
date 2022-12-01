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

	Router diface.IRouter
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

	fmt.Printf("[Server] %s, Version %s is serving...\n", s.Name, s.Version)
	select {}
}

// AddRouter ...
func (s *server) AddRouter(router diface.IRouter) {
	s.Router = router
	fmt.Println("[Server] add router")
}
