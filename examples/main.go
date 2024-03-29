package main

import (
	"fmt"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/dnet"
)

type router1 struct {
	dnet.BaseRouter
}

func (r *router1) Handle(req diface.IRequest) {
	fmt.Printf("[Router] custom hanle method 1 for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	fmt.Printf("[Router] receive msg, ID = %d, dataLen = %d, data:%s\n", req.GetMessageID(), req.GetDataLength(), req.GetData())

	if err := req.GetConnection().SendMsg(1, []byte("Hello, This is response for msg1")); err != nil {
		fmt.Printf("[Router] failed to send msg back to client for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	}
}

type router2 struct {
	dnet.BaseRouter
}

func (r *router2) Handle(req diface.IRequest) {
	fmt.Printf("[Router] custom hanle method 2 for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	fmt.Printf("[Router] receive msg, ID = %d, dataLen = %d, data:%s\n", req.GetMessageID(), req.GetDataLength(), req.GetData())

	if err := req.GetConnection().SendMsg(2, []byte("Hello, This is response for msg2")); err != nil {
		fmt.Printf("[Router] failed to send msg back to client for connection, ID = %d\n", req.GetConnection().GetConnectionID())
	}
}

func main() {
	srv := dnet.NewServer()

	srv.SetOnConnStart(func(conn diface.IConnection) {
		fmt.Printf("[Hook] connectionID = %d start\n", conn.GetConnectionID())
		conn.SetProperty("key", "custom key value")
	})
	srv.SetOnConnClose(func(conn diface.IConnection) {
		fmt.Printf("[Hook] connectionID =%d close\n", conn.GetConnectionID())
		value, _ := conn.GetProperty("key")
		fmt.Println("value:", value)
	})

	srv.AddRouter(1, &router1{})
	srv.AddRouter(2, &router2{})

	srv.Serve()
}
