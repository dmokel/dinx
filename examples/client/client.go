package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var exit = make(chan bool)

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("[Client] dial tcp error: ", err)
		exit <- true
		return
	}

	go reader(conn)
	go writer(conn)
}

func reader(conn net.Conn) {
	defer close(exit)

	pack := newPack()
	for {
		headBuf := make([]byte, pack.GetHeadLen())
		if _, err := io.ReadFull(conn, headBuf); err != nil {
			fmt.Println("failed to read head buf from server, err:", err)
			break
		}
		msg, err := pack.Unpack(headBuf)
		if err != nil {
			fmt.Println("failed to upack head buf, err:", err)
			break
		}
		dataBuf := make([]byte, msg.dataLen)
		if msg.dataLen > 0 {
			if _, err := io.ReadFull(conn, dataBuf); err != nil {
				fmt.Println("failed to read data buf from server, err:", err)
				break
			}
		}
		msg.data = dataBuf
		fmt.Printf("receive msg from server, msgID = %d, dataLen = %d, data:%s\n", msg.msgID, msg.dataLen, msg.data)
	}
}

func writer(conn net.Conn) {
	pack := newPack()
	for {
		select {
		case <-exit:
			return
		default:
			if err := sendMsg1(conn, pack); err != nil {
				fmt.Println("send msg1 err:", err)
			}
			if err := sendMsg2(conn, pack); err != nil {
				fmt.Println("send msg 2 err:", err)
			}

			time.Sleep(3 * time.Second)
		}
	}
}

func sendMsg1(conn net.Conn, pack *pack) error {
	msg := &message{}
	msg.msgID = 1
	dataBuf := []byte("Hello, This is client one")
	msg.data = dataBuf
	msg.dataLen = uint32(len(dataBuf))
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

func sendMsg2(conn net.Conn, pack *pack) error {
	msg := &message{}
	msg.msgID = 2
	dataBuf := []byte("Hello, This is client two")
	msg.data = dataBuf
	msg.dataLen = uint32(len(dataBuf))
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
	go client()
	<-exit
}
