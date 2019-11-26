package ziface

type IHook interface {
	// 设置该Server的连接创建时Hook函数
	SetOneConnStart(func(conn IConnection))
	// 设置该Server的连接断开时的Hook函数
	SetOneConnStop(func(conn IConnection))
	// 调用连接OnConnStart Hook函数
	CallOneConnStart(conn IConnection)
	// 调用连接OnConnStop Hook函数
	CallOneConnStop(conn IConnection)
}
