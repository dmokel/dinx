package main

import (
	"fmt"
	"io"
	"net"
	"time"

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
		msg := dnet.NewMessage()
		msg.SetMsgID(0)
		dataBuf := []byte("Hello, This is client one")
		msg.SetData(dataBuf)
		msg.SetDataLen(uint32(len(dataBuf)))
		buf, err := pack.Pack(msg)
		if err != nil {
			fmt.Println("failed to pack msg, err:", err)
			break
		}
		if _, err := conn.Write(buf); err != nil {
			fmt.Println("failed to write buffer conn, err:", err)
			break
		}

		headBuf := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(conn, headBuf); err != nil {
			fmt.Println("failed to read head buf from server, err:", err)
			break
		}
		msg, err = pack.Unpack(headBuf)
		if err != nil {
			fmt.Println("failed to upack head buf, err:", err)
			break
		}
		dataBuf = make([]byte, msg.GetDataLen())
		if msg.GetDataLen() > 0 {
			if _, err := io.ReadFull(conn, dataBuf); err != nil {
				fmt.Println("failed to read data buf from server, err:", err)
				break
			}
		}
		msg.SetData(dataBuf)
		fmt.Printf("receive msg from server, msgID = %d, dataLen = %d, data:%s\n", msg.GetMsgID(), msg.GetDataLen(), msg.GetData())
		time.Sleep(2 * time.Second)
	}
	exit <- true
}

func main() {
	exit := make(chan bool)
	go client(exit)
	<-exit
}
