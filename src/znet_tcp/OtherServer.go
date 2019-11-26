package znet_tcp

import (
	"fmt"
	"time"
	"ziface"
	"utils"
	"net"
	"message"
	"LuaState"
)

//iServer 接口实现，定义一个Server服务类
type OtherServer struct {
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
	ownServer 	ziface.IServer
}

/*
  创建一个服务器句柄
 */
func NewOtherServer (name string) ziface.IServer {
	s:= &OtherServer {
		Name :utils.GlobalServer.ServerName,
		IPVersion:"tcp",
		Host:"--",
		ConnMgr:NewConnMgr(),
	}
	s.msgHandle = NewMsgHandle(s)

	return s
}
//============== 实现 ziface.IServer 里的全部接口方法 ========

func (s *OtherServer) initConnection(dialConn ziface.IConnection, other_conf utils.OtherServerConfig){
	// 添加连接
	s.ConnMgr.Add(dialConn)

	go dialConn.Start()

	// 确认服务器Id
	dialConn.SendBufMsg(message.NewServerIdMsg(utils.GlobalServer.Id, utils.GlobalServer.Name, utils.GlobalServer.ServerName), false)

	go LuaState.AddServer(other_conf.Id, other_conf.Name, other_conf.ServerName)
}

// 链接其他服务
func (s *OtherServer) ConnectionOther() {
	count := 0
	waitConnList := map[utils.OtherServerConfig]bool{}

	for _,v := range utils.OtherServerConfigList {
		tcp_addr,_ := net.ResolveTCPAddr("tcp", v.Host)
		conn, err := net.DialTCP("tcp", nil, tcp_addr)
		if err == nil {
			dialConn := NewConntion(s, conn, v.Id, s.msgHandle)
			fmt.Println("v.Host:", v.Host)

			s.initConnection(dialConn, v)
		} else {
			fmt.Println("err:", err)
			waitConnList[v] = true
			count++
		}
	}

	fmt.Println(waitConnList)
	for ;count > 0; {
		for v,_ := range waitConnList {
			tcp_addr,_ := net.ResolveTCPAddr("tcp", v.Host)
			conn, err := net.DialTCP("tcp", nil, tcp_addr)
			if err == nil {
				dialConn := NewConntion(s, conn, v.Id, s.msgHandle)

				s.initConnection(dialConn, v)

				count--

				delete(waitConnList, v)
			} else {
				//log.Fatalln("dailing error: ", err)
				fmt.Println("dailing error: ", err)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

//开启网络服务
func (s *OtherServer) Start() {
	fmt.Printf("[START] Server name: %s,listenner at Host: %s is starting\n", s.Name, s.Host)
	fmt.Printf("[Zinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		utils.GlobalServer.Version,
		utils.GlobalServer.MaxConn,
		utils.GlobalServer.MaxPacketSize)
	fmt.Printf("[Worker] WorkerPoolSize: %d, MaxWorkerTaskLen: %d\n",
		utils.GlobalServer.WorkerPoolSize,
		utils.GlobalServer.MaxWorkerTaskLen)

	//1 启动worker工作池机制
	s.msgHandle.StartWorkerPool()

	//3 开启一个go去做服务端Linster业务
	go s.ConnectionOther()
}

func (s *OtherServer)GetConnMgr() ziface.IConnManager{
	return  s.ConnMgr
}

func (s *OtherServer) Stop() {
	fmt.Println("[STOP] Zinx server , name " , s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

func (s *OtherServer) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出

	for {
		time.Sleep(10*time.Second)
	}
}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *OtherServer)AddRouter(msgId uint32, router ziface.IRouter) {
	s.msgHandle.AddRouter(msgId, router)
}

func (s *OtherServer)GetHook() ziface.IHook {
	return s.GetIHook()
}

func (s *OtherServer)SetServer(ser ziface.IServer) {
	s.ownServer = ser
}

func (s *OtherServer)GetServer() ziface.IServer {
	return s.ownServer
}

func (s *OtherServer)SendTo(id uint32, msg message.IMessage) {
	conn, err := s.ConnMgr.GetById(id)
	if err == nil {
		conn.SendBufMsg(msg, false)
	} else {
		if s.ownServer != nil {
			conn, err = s.ownServer.GetConnMgr().GetById(id)
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

