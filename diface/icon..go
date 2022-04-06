package diface

import (
	"context"
	"net"
)

// IConn ...
type IConn interface {
	Start()                   // 启动连接, 让当前连接开始工作
	Stop()                    // 停止连接, 结束当前连接状态
	Context() context.Context // 返回ctx,

	GetTcpConn() *net.TCPConn // 从连接中获取原始tcp conn/socket conn
	GetConnID() uint32        // 获取连接ID
	RemoteAddr() net.Addr     // 获取对端地址信息

	SendMsg(msgID uint32, data []byte) error     // 将消息发送给对端（无缓冲）
	SendBuffMsg(msgID uint32, data []byte) error // 将消息发送给对端（有缓冲）

	setAttribute(key string, value interface{})   // 设置连接属性
	GetAttribute(key string) (interface{}, error) // 获取连接属性
	RemoveAttribute(key string)                   // 移除连接属性
}
