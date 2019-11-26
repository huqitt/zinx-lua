package Router

import (
	"ziface"
	"github.com/yuin/gopher-lua"
	"message"
)

type ServerIdRouter struct {
	ziface.BaseRouter
}

func (r *ServerIdRouter) PreHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PreHandle:", message)
}

func (r *ServerIdRouter) Handle(conn ziface.IConnection, lState *lua.LState, c_msg message.IMessage) {
	i_msg := c_msg.(interface{})
	switch msg := i_msg.(type) {
	case *message.ServerIdMsg:
		conn.GetServer().GetConnMgr().SetConnId(conn.GetConnID(), msg.ServerId)
	}
}

func (r *ServerIdRouter) PostHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PostHandle:", message)
}
