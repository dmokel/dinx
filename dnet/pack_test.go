package dnet

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestPack(t *testing.T) {
	exit := make(chan bool)

	listener, err := net.Listen("tcp", "127.0.0.1:8999")
	if err != nil {
		t.Fatal(err)
	}

	go func() {

		conn, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}

		pack := newPack()
		msgCount := 0
		for {
			headBuf := make([]byte, pack.GetHeadLen())
			_, err = io.ReadFull(conn, headBuf)
			if err != nil {
				t.Fatal(err)
			}

			msg, err := pack.Unpack(headBuf)
			if err != nil {
				t.Fatal(err)
			}

			if msg.GetDataLen() > 0 {
				dataBuf := make([]byte, msg.GetDataLen())
				_, err = io.ReadFull(conn, dataBuf)
				if err != nil {
					t.Fatal(err)
				}
				msg.SetData(dataBuf)
				fmt.Printf("receive msg, msgID=%d, dataLen=%d, data=%s\n", msg.GetMsgID(), msg.GetDataLen(), msg.GetData())
				msgCount++
			}

			if msgCount >= 2 {
				exit <- true
			}
		}
	}()

	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8999")

	data1 := []byte("Hello,World. This is Mokel.")
	msg1 := &message{
		msgID:   2,
		dataLen: uint32(len(data1)),
		data:    data1,
	}

	data2 := []byte("HelloWorld")
	msg2 := &message{
		msgID:   1,
		dataLen: uint32(len(data2)),
		data:    data2,
	}

	pack := newPack()
	buf1, err := pack.Pack(msg1)
	if err != nil {
		t.Fatal(err)
	}
	buf2, err := pack.Pack(msg2)

	buf := make([]byte, 0)
	buf = append(buf, buf1...)
	buf = append(buf, buf2...)
	_, err = conn.Write(buf)
	if err != nil {
		t.Fatal(err)
	}

	<-exit
}
