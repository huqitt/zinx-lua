package znet_tcp

import (
	"znet_websocket"
	"ziface"
	"fmt"
)

func onConnStartCall(conn ziface.IConnection){
	fmt.Println("start:", conn.RemoteAddr())
}

func onConnStopCall(conn ziface.IConnection){
	fmt.Println("stop:", conn.RemoteAddr())
}

func main() {
	serve := NewServer("test_serve")

	serve.GetHook().SetOneConnStop(onConnStopCall)
	serve.GetHook().SetOneConnStart(onConnStartCall)

	r := znet_websocket.Router{}

	serve.AddRouter(0, &r)

	serve.Serve()

	otherServer := NewOtherServer("test_other_server")
	otherServer.GetHook().SetOneConnStop(onConnStopCall)
	otherServer.GetHook().SetOneConnStart(onConnStartCall)
	otherServer.AddRouter(0, &r)

	otherServer.Serve()
}
