package ziface

import "message"

type IServer interface {
	// 启动服务器方法
	Start()
	// 停止服务器方法
	Stop()
	// 开启服务器方法
	Serve()
	// 添加路由
	AddRouter(msgId uint32, router IRouter)
	// 获取链接管理器
	GetConnMgr() IConnManager
	//
	GetHook() IHook
	//
	SetServer(ser IServer)
	//
	GetServer() IServer
	//
	SendTo(id uint32, msg message.IMessage)
}
