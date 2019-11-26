package znet_tcp

import (
	"ziface"
	"sync"
	"fmt"
	"errors"
	"LuaState"
	"utils"
)
/*
    连接管理模块
*/
type ConnManager struct {
	connections 	map[uint32]ziface.IConnection //管理的连接信息
	connectionsKey 	map[uint32]ziface.IConnection //管理的连接信息
	connLock    	sync.RWMutex                  //读写连接的读写锁
	connKeyLock    	sync.RWMutex                  //读写连接的读写锁
}

/*
    创建一个链接管理
*/
func NewConnMgr() *ConnManager {
	return &ConnManager{
		connections:make(map[uint32] ziface.IConnection),
		connectionsKey:make(map[uint32] ziface.IConnection),
	}
}
//添加链接
func (mgr *ConnManager) Add(conn ziface.IConnection) {
	mgr.connLock.Lock()

	defer mgr.connLock.Unlock()
	//将conn连接添加到ConnMananger中
	mgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", mgr.Len())
}

// 添加链接
func (mgr *ConnManager) SetConnId(connID uint32, id uint32) {
	mgr.connKeyLock.Lock()

	defer mgr.connKeyLock.Unlock()
	//
	conn,ok := mgr.connections[connID]
	if ok {
		mgr.connectionsKey[id] = conn
	}

	//fmt.Println("connection add to ConnManager successfully: conn num = ", mgr.Len())
}
// 根据id获取连接
func (mgr *ConnManager) GetById(id uint32) (ziface.IConnection, error){
	mgr.connKeyLock.Lock()

	defer mgr.connKeyLock.Unlock()

	conn,ok := mgr.connectionsKey[id]
	if ok {
		return conn, nil
	} else {
		return nil, errors.New("不存在id:" + utils.ToString(id))
	}
}
//获取当前连接数量
func (mgr *ConnManager)Len()int {
	return len(mgr.connections)
}
//删除连接
func (mgr *ConnManager)Remove(conn ziface.IConnection) {
	if _,ok := mgr.connections[conn.GetConnID()]; ok {
		fmt.Println("start Remove ConnID =", conn.GetConnID())

		mgr.connLock.Lock()

		defer mgr.connLock.Unlock()

		delete(mgr.connections, conn.GetConnID())

		mgr.connKeyLock.Lock()

		defer mgr.connKeyLock.Unlock()

		for k,v := range mgr.connectionsKey {
			if v == conn {
				delete(mgr.connectionsKey, k)
				break
			}
		}

		fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", mgr.Len())
	}
}
//清除并停止所有连接
func (mgr *ConnManager)ClearConn() {
	mgr.connLock.Lock()

	defer mgr.connLock.Unlock()

	for k,v := range mgr.connections {

		delete(mgr.connections, k)

		v.Stop()
	}

	mgr.connKeyLock.Lock()

	defer mgr.connKeyLock.Unlock()

	for k,_ := range mgr.connectionsKey {
		delete(mgr.connectionsKey, k)
	}

	// 关闭所有虚拟机
	LuaState.ClearLStateAll()

	fmt.Println("Clear All Connections successfully: conn num = ", mgr.Len())
}
//利用ConnID获取链接
func (mgr *ConnManager)Get(connID uint32) (ziface.IConnection, error){

	mgr.connLock.Lock()

	defer mgr.connLock.Unlock()

	conn,ok := mgr.connections[connID]
	if ok {
		return conn, nil
	}else{
		return nil, errors.New("connection not found")
	}
}


