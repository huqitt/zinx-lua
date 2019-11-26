package znet_tcp

import (
	"ziface"
	"github.com/yuin/gopher-lua"
)

type Request struct {
	data 		[]byte
	Count		int
	conn 		ziface.IConnection
	lState		*lua.LState
}

func (r *Request) GetConnection() ziface.IConnection{
	return r.conn
}

func (r *Request) GetData() []byte{
	return r.data
}

func (r *Request) GetCount() int{
	return r.Count
}

func (r *Request) SetLState(l *lua.LState){
	r.lState = l
}

func (r *Request) GetLState() *lua.LState{
	return r.lState
}
