package ziface

import "github.com/yuin/gopher-lua"

/*
    IRequest 接口：
    实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface{
	GetConnection() IConnection //获取请求连接信息
	GetData() []byte            //获取请求消息的数据
	GetCount() int            	//获取请求消息的长度
	SetLState(l *lua.LState)	//设置lua虚拟机
	GetLState() *lua.LState		//获取lua虚拟机
}