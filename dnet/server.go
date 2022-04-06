package dnet

import "github.com/dmokel/dinx/diface"

// Server ...
type Server struct {
	Name    string
	Network string
	IP      string
	Port    int

	connMgr    diface.IConnMgr    // conn pool mgr module
	msgHandler diface.IMsgHandler // msg handler module

	OnConnOpen  func(conn diface.IConn)
	OnConnClose func(conn diface.IConn)
}

// NewServer ...
func NewServer(opts ...Option) diface.IServer {
	s := &Server{
		Name:    "default",
		Network: "tcp",
		IP:      "127.0.0.1",
		Port:    8888,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start ...
func (s *Server) Start() {}

// Stop ...
func (s *Server) Stop() {}

// Serve ...
func (s *Server) Serve() {}
