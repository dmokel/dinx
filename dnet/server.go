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
}

var _ diface.IServer = &server{}

// NewServer used to create the server instance
func NewServer() diface.IServer {
	return &server{
		Name:    "default",
		Network: "tcp",
		IP:      "127.0.0.1",
		Port:    8999,
	}
}

func (s *server) Start() {
	go func() {
		localAddr, err := net.ResolveTCPAddr(s.Network, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Start] failed to resolve tcp addr, err: ", err)
			return
		}

		listener, err := net.ListenTCP(s.Network, localAddr)
		if err != nil {
			fmt.Println("[Start] failed to create tcp listener, err: ", err)
			return
		}
		fmt.Printf("[Start] success listening tcp connection at %s:%d\n", s.IP, s.Port)

		for {
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("[Listening] failed to accept tcp connection, err: ", err)
				continue
			}

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := tcpConn.Read(buf)
					if err != nil {
						fmt.Println("read buf error: ", err)
						continue
					}

					_, err = tcpConn.Write(buf[:cnt])
					if err != nil {
						fmt.Println("write buf error: ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *server) Stop() {}

func (s *server) Serve() {
	s.Start()

	fmt.Println("[Serve] Dinx Server Serve...")
	select {}
}
