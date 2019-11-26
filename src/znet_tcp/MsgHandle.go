package znet_tcp

import (
	"ziface"
	"message"
	"fmt"
	"strconv"
	"utils"
	"github.com/yuin/gopher-lua"
	"LuaState"
)

type MsgHandle struct {
	Apis 			map[uint32]ziface.IRouter //
	WorkPoolSize	uint32
	TaskQueue		[]chan	ziface.IRequest
	LuaState		*lua.LState
	LStateId		int
	Server			ziface.IServer
}

func NewMsgHandle(ser ziface.IServer) *MsgHandle{
	return &MsgHandle{
		Apis:make(map[uint32]ziface.IRouter, 10),
		WorkPoolSize:utils.GlobalServer.WorkerPoolSize,
		TaskQueue:make([]chan ziface.IRequest, utils.GlobalServer.WorkerPoolSize),
		Server:ser,
	}
}

func (m *MsgHandle) InitLState(L *lua.LState){
	SendMsgFunc := func(L *lua.LState) int {
		ser_id := L.ToInt(1)
		L.Pop(1)
		table := L.ToString(2)
		L.Pop(2)
		funcName := L.ToString(3)
		L.Pop(3)
		param := L.ToString(4)
		L.Pop(4)
		msg := message.NewLuaMessage(table, funcName, param)
		m.Server.SendTo(uint32(ser_id), msg)
		return 1
	}

	PanicFunc := func(L *lua.LState) int {
		str := L.ToString(1)
		L.Pop(1)
		panic(str)
		return  1
	}

	mainTable := L.GetGlobal("Main").(*lua.LTable)
	mainTable.RawSetString("SendMsg", L.NewFunction(SendMsgFunc))
	mainTable.RawSetString("Panic", L.NewFunction(PanicFunc))
}

//Worker工作流程
func (m *MsgHandle) OneWorker(workerID int, ch chan ziface.IRequest){
	// 启动lua虚拟机
	LState, LStateId := LuaState.GetLState()
	// 初始化虚拟机
	m.InitLState(LState)
	//fmt.Println("Worker ID = ", workerID, " is started.")
	for{
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case data := <- ch:
			// 设置消息对应的虚拟机
			data.SetLState(LState)
			// 处理消息
			m.DoMsgHandler(data)
		}
	}
	// 归还虚拟机
	LuaState.CloseLState(LStateId)
}
//将消息交给TaskQueue,由worker进行处理
func (m *MsgHandle) SendToTaskQueue(request ziface.IRequest){
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % m.WorkPoolSize
	//将请求消息发送给任务队列
	m.TaskQueue[workerID] <- request
}
//启动worker工作池
func (m *MsgHandle)StartWorkerPool(){
	for i := 0; i < int(m.WorkPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalServer.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go m.OneWorker(i, m.TaskQueue[i])
	}
}

//马上以非阻塞方式处理消息
func (m *MsgHandle) DoMsgHandler(request ziface.IRequest){
	defer fmt.Println(request.GetData()[:request.GetCount()])

	propertyList := message.GetPropertyList(request.GetData())
	msg := message.GetMessage(propertyList)
	msgId := message.GetMsgId(propertyList)
	handler,ok := m.Apis[msgId]
	if !ok {
		fmt.Println("api msgId = ", msgId, " is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request.GetConnection(), request.GetLState(), msg)
	handler.Handle(request.GetConnection(), request.GetLState(), msg)
	handler.PostHandle(request.GetConnection(), request.GetLState(), msg)
}

//为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	m.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (m *MsgHandle) GetServer() ziface.IServer{
	return m.Server
}
