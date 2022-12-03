package main

import (
	"fmt"
	"net"
	"time"

	"github.com/dmokel/dinx/diface"
	"github.com/dmokel/dinx/dnet"
)

func client(exit chan<- bool) {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("[Client] dial tcp error: ", err)
		exit <- true
		return
	}
	pack := dnet.NewPack()
	for {
		if err := sendMsg1(conn, pack); err != nil {
			break
		}
		if err := sendMsg2(conn, pack); err != nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	exit <- true
}

func sendMsg1(conn net.Conn, pack diface.IPack) error {
	msg := dnet.NewMessage()
	msg.SetMsgID(1)
	dataBuf := []byte("Hello, This is client one")
	msg.SetData(dataBuf)
	msg.SetDataLen(uint32(len(dataBuf)))
	buf, err := pack.Pack(msg)
	if err != nil {
		fmt.Println("failed to pack msg, err:", err)
		return err
	}
	if _, err := conn.Write(buf); err != nil {
		fmt.Println("failed to write buffer conn, err:", err)
		return err
	}
	return nil
}

func sendMsg2(conn net.Conn, pack diface.IPack) error {
	msg := dnet.NewMessage()
	msg.SetMsgID(2)
	dataBuf := []byte("Hello, This is client two")
	msg.SetData(dataBuf)
	msg.SetDataLen(uint32(len(dataBuf)))
	buf, err := pack.Pack(msg)
	if err != nil {
		fmt.Println("failed to pack msg, err:", err)
		return err
	}
	if _, err := conn.Write(buf); err != nil {
		fmt.Println("failed to write buffer conn, err:", err)
		return err
	}
	return nil
}

func main() {
	exit := make(chan bool)
	go client(exit)
	<-exit
}
