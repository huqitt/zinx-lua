package znet_tcp

import "ziface"

type Hook struct {
	//
	connStartCall(func(conn ziface.IConnection))
	//
	connStopCall(func(conn ziface.IConnection))
}

func (h *Hook) GetIHook() ziface.IHook{
	return h
}

//============== 实现 ziface.IHook 里的全部接口方法 ==============
func(h *Hook)SetOneConnStart(function func(conn ziface.IConnection)){
	h.connStartCall = function
}

func(h *Hook)SetOneConnStop(function func(conn ziface.IConnection)){
	h.connStopCall = function
}

func(h *Hook)CallOneConnStart(conn ziface.IConnection){
	if h.connStartCall != nil {
		h.connStartCall(conn)
	}
}

func(h *Hook)CallOneConnStop(conn ziface.IConnection){
	if h.connStopCall != nil {
		h.connStopCall(conn)
	}
}