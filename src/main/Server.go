package main

import (
	"ziface"
	"fmt"
	"znet_tcp"
	"message"
	"Router"
)

func onConnStartCall(conn ziface.IConnection){
	fmt.Println("start:", conn.RemoteAddr())
}

func onConnStopCall(conn ziface.IConnection){
	fmt.Println("stop:", conn.RemoteAddr())
}

func main() {
	serve := znet_tcp.NewServer("test_serve")

	serve.GetHook().SetOneConnStop(onConnStopCall)
	serve.GetHook().SetOneConnStart(onConnStartCall)

	heart := &Router.HeartRouter{}
	serve.AddRouter(message.MSG_HEART, heart)

	lua := &Router.LuaMsgRouter{}
	serve.AddRouter(message.MSG_LUA_STRING, lua)

	serve.Serve()

	otherServer := znet_tcp.NewOtherServer("test_other_server")
	otherServer.GetHook().SetOneConnStop(onConnStopCall)
	otherServer.GetHook().SetOneConnStart(onConnStartCall)
	otherServer.AddRouter(message.MSG_HEART, heart)

	otherServer.SetServer(serve)
	serve.SetServer(otherServer)

	otherServer.Serve()
}
