package Router

import (
	"ziface"
	"message"
	"github.com/yuin/gopher-lua"
	"fmt"
	"reflect"
)

type HeartRouter struct {
	ziface.BaseRouter
}

func (r *HeartRouter) PreHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PreHandle:", message)
}

func (r *HeartRouter) Handle(conn ziface.IConnection, lState *lua.LState, c_msg message.IMessage) {
	//fmt.Println("\nCall Router Handle:", message)
	i_msg := c_msg.(interface{})
	conn.SendBufMsg(c_msg, false)
	fmt.Println("i_msg.type =", reflect.TypeOf(i_msg))
	switch i_msg.(type) {
	case *message.Heart:
		//fmt.Println("\nCall Router Handle:", message)
		err := lState.CallByParam(
			lua.P{
				Fn: lState.GetGlobal("Main").(*lua.LTable).RawGetString("OnMessage").(*lua.LFunction),
				NRet: 0,
				Protect: true,
			},
			lua.LString("Main"),
			lua.LString("OnMessage"),
			lua.LString("{1,2,3}"),
		)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("call -> Main.OnMessage")
		}
	}
}

func (r *HeartRouter) PostHandle(conn ziface.IConnection, lState *lua.LState, message message.IMessage) {
	//fmt.Println("\nCall Router PostHandle:", message)
}