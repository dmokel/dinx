package main

import (
	"fmt"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/dnet"
)

type router struct {
	dnet.BaseRouter
}

func (r *router) Handle(req diface.IRequest) {
	fmt.Printf("[Router] custom hanle method for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	fmt.Printf("[Router] receive msg, ID = %d, dataLen = %d, data:%s\n", req.GetMessageID(), req.GetDataLength(), req.GetData())

	if err := req.GetConnection().SendMsg(1, []byte("Hello, This is Server")); err != nil {
		fmt.Printf("[Router] failed to send msg back to client for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	}
}

func main() {
	srv := dnet.NewServer()

	srv.AddRouter(&router{})

	srv.Serve()
}
