package Router

import (
	"ziface"
	"message"
	"github.com/yuin/gopher-lua"
)

type LuaMsgRouter struct {
	ziface.BaseRouter
}

func (r *LuaMsgRouter) PreHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PreHandle:", message)
}

func (r *LuaMsgRouter) Handle(conn ziface.IConnection, lState *lua.LState, c_msg message.IMessage) {
	i_msg := c_msg.(interface{})
	switch msg := i_msg.(type) {
	case *message.LuaMessage:
		//fmt.Println("\nCall Router Handle:", message)
		err := lState.CallByParam(
			lua.P{
				Fn: lState.GetGlobal("Main").(*lua.LTable).RawGetString("OnMessage").(*lua.LFunction),
				NRet: 0,
				Protect: true,
			},
			lua.LString(msg.Table),
			lua.LString(msg.Function),
			lua.LString(msg.Param),
		)
		if err != nil {
			panic(err)
		}
	}
}

func (r *LuaMsgRouter) PostHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PostHandle:", message)
}