package znet_tcp

import (
	"fmt"
	"net"
	"ziface"
	"utils"
	"message"
)

//iServer 接口实现，定义一个Server服务类
type Server struct {
	// hook
	Hook
	//服务器的名称
	Name string
	//tcp4 or other
	IPVersion string
	//服务绑定的IP地址
	Host string
	//服务绑定的端口
	Port int
	//当前Server由用户绑定的回调router,也就是Server注册的链接对应的处理业务
	msgHandle	ziface.IMsgHandle
	// 链接管理器
	ConnMgr		ziface.IConnManager
	//
	otherServer	ziface.IServer
}


/*
  创建一个服务器句柄
 */
func NewServer (name string) ziface.IServer {
	s := &Server {
		Name :utils.GlobalServer.Name,
		IPVersion:"tcp",
		Host:utils.GlobalServer.Host,
		ConnMgr:NewConnMgr(),
	}
	s.msgHandle = NewMsgHandle(s)

	return s
}

//============== 实现 ziface.IServer 里的全部接口方法 ========

// 监听客户端链接
func (s *Server) ListenerClient() {
	//0 启动worker工作池机制
	s.msgHandle.StartWorkerPool()

	//1 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(s.IPVersion, s.Host)
	if err != nil {
		fmt.Println("resolve tcp addr err: ", err)
		return
	}

	//2 监听服务器地址
	listenner, err:= net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("listen", s.IPVersion, "err", err)
		return
	}

	//已经监听成功
	fmt.Println("start Zinx-Lua server  ", s.Name, " succ, now listenning...")

	//TODO server.go 应该有一个自动生成ID的方法
	var cid uint32
	cid = 0

	//3 启动server网络连接业务
	for {
		//3.1 阻塞等待客户端建立连接请求
		conn, err := listenner.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err ", err)
			continue
		}
		fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

		//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
		if s.ConnMgr.Len() >= utils.GlobalServer.MaxConn {
			conn.Close()
			continue
		}

		//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		dealConn := NewConntion(s, conn, cid, s.msgHandle)

		//3.4 添加连接
		s.ConnMgr.Add(dealConn)
		//3.5 设置客户端 唯一标识符
		s.ConnMgr.SetConnId(cid, cid)

		cid ++

		//3.4 启动当前链接的处理业务
		go dealConn.Start()
	}
}

//开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at Host: %s is starting\n", s.Name, s.Host)
	fmt.Printf("[Zinx-Lua] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalServer.Version,
		utils.GlobalServer.MaxConn,
		utils.GlobalServer.MaxPacketSize)
	fmt.Printf("[Worker] WorkerPoolSize: %d, MaxWorkerTaskLen: %d\n",
		utils.GlobalServer.WorkerPoolSize,
		utils.GlobalServer.MaxWorkerTaskLen)

	//1 启动worker工作池机制
	s.msgHandle.StartWorkerPool()

	//2 开启一个go去做服务端Linster业务
	go s.ListenerClient()
}

//停止服务
func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name " , s.Name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

//运行服务
func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	//select{}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server)AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandle.AddRouter(msgId, router)
}

//得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func (s *Server)GetHook() ziface.IHook {
	return s.GetIHook()
}

func (s *Server)SetServer(ser ziface.IServer) {
	s.otherServer = ser
}

func (s *Server)GetServer() ziface.IServer {
	return s.otherServer
}

func (s *Server)SendTo(id uint32, msg message.IMessage) {
	conn, err := s.ConnMgr.GetById(id)
	if err == nil {
		conn.SendBufMsg(msg, false)
	} else {
		if s.otherServer != nil {
			conn, err = s.otherServer.GetConnMgr().GetById(id)
			if err == nil {
				conn.SendBufMsg(msg, false)
			} else {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}





