package main

import (
	"fmt"
	"net"
	"time"

	"github.com/dmokel/dinx/dnet"
)

func main() {
	go server()

	time.Sleep(3 * time.Second)

	go client()

	select {}
}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("[Client] dial tcp error: ", err)
		return
	}
	for {
		cnt, err := conn.Write([]byte("HelloWorld"))
		if err != nil {
			fmt.Println("[Client] write bytes error: ", err)
			continue
		}
		fmt.Println("[Client] send msg: HelloWorld")

		buf := make([]byte, cnt)
		_, err = conn.Read(buf)
		if err != nil {
			fmt.Println("[Client] read bytes error: ", err)
			continue
		}
		fmt.Println("[Client] receive msg:", string(buf))

		time.Sleep(2 * time.Second)
	}
}

func server() {
	srv := dnet.NewServer()
	srv.Serve()
}
