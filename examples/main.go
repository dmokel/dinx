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
	fmt.Println("[Router] custom hanle method")
	if err := req.GetConnection().Send(req.GetData()[:req.GetDataLength()]); err != nil {
		fmt.Println("[Router] failed to send data back to client, err:", err)
	}
}

func main() {
	srv := dnet.NewServer()

	srv.AddRouter(&router{})

	srv.Serve()
}
