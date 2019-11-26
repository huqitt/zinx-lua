package ziface

import (
	"net"
	"message"
)

type IConnection interface {
	// 启动链接，让当前链接开始工作
	Start()
	// 停止链接，结束当前链接状态
	Stop()
	// 获取Server
	GetServer() IServer
	// 获取当前websocket的Con
	GetConnection() net.Conn
	// 获取当前链接ID
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 发送给客户端消息
	SendMsg(data message.IMessage, retrans bool) error
	// 发送给客户端消息
	SendBufMsg(data message.IMessage, retrans bool) error
	// 获取属性
	GetProperty() IProperty
}