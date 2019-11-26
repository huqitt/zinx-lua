package znet_tcp

import (
	"fmt"
	"net"
	"ziface"
	"errors"
	"message"
	"utils"
)

const (
	UINT16_MAx = 65536
)

type Connection struct {
	// 继承Property
	Property
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	// Server
	Server ziface.IServer
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//该连接的处理方法router
	msgHandle  ziface.IMsgHandle

	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool

	// 发送给客户端的管道，无缓存
	msgChan chan *DataPack
	// 发送给客户端的管道，有缓存
	msgBufChan chan *DataPack
	// 发给客户端的消息序号
	_seq uint16
	// 客户端发来的消息确认号
	_ack uint16
	// 可靠消息列表
	sendList []*DataPack
}


//创建连接的方法
func NewConntion(server ziface.IServer, conn *net.TCPConn, connID uint32, hander  ziface.IMsgHandle) *Connection{
	c := &Connection{
		Server:server,
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		msgHandle: hander,
		ExitBuffChan: make(chan bool, 1),
		msgChan:make(chan *DataPack),
		msgBufChan:make(chan *DataPack, utils.GlobalServer.MaxMsgChanLen),
		Property:Property{
			property:make(map[string]interface{}),
		},
		_seq : 0,
		_ack : 0,
		sendList:make([]*DataPack, 0),
	}

	return c
}

func (c *Connection)GetServer() ziface.IServer{
	return c.Server
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		//读取我们最大的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("[%d]recv buf err:%s", c.ConnID, err)
			c.ExitBuffChan <- true
			continue
		}
		if cnt > 0 {
			fmt.Println("[%d]收到数据包大小：%d", c.ConnID, cnt)

			fmt.Println("收到的数据包：", buf)

			buf = buf[0:cnt]

			fmt.Println("切片后的数据包：", buf)
			// 解析拿到数据包
			datapack := NewUnDataPack(buf)
			fmt.Println("解包后的数据包：", datapack.Body)
			// 设置客户端消息消息确认号为本次消息序号
			c._ack = datapack.Seq
			// 根据消息号移除需要确认的消息
			for sl_len := len(c.sendList); sl_len > 0; {
				top := c.sendList[sl_len - 1]
				if datapack.Ack - top.Seq <= UINT16_MAx / 2 {
					if datapack.Ack < top.Seq {
						break
					}
				} else {
					if datapack.Ack > top.Seq {
						break
					}
				}
				c.sendList = c.sendList[:sl_len - 1]
				sl_len = len(c.sendList)
			}

			//得到当前客户端请求的Request数据
			req := Request{
				conn:c,
				data:datapack.Body,
				Count:cnt - datapack.GetHeadLen(),
			}
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			c.msgHandle.SendToTaskQueue(&req)
		}
	}
}

func (c *Connection)StartWriter(){
	for{
		select {
		case pack := <-c.msgChan:
			// 发送给客户端数据
			_,err := c.Conn.Write(pack.Data)
			if err != nil {
				fmt.Println(c.Conn.RemoteAddr().String(), " - Send Data error:", err, ".Conn Writer exit")
				return
			}
			fmt.Printf("[%d]发送消息长度：%d\n", c.ConnID,len(pack.Body))
		case pack, ok := <-c.msgBufChan:
			if ok {
				// 发送给客户端数据
				_,err := c.Conn.Write(pack.Data)
				fmt.Println("发送：", pack.Data)
				if err != nil {
					fmt.Println(c.Conn.RemoteAddr().String(), " - Send Data error:", err, ".Conn Writer exit")
					return
				}
			} else {
				fmt.Println(c.Conn.RemoteAddr().String(), " - closed! Conn Writer exit")
				return
			}
		case <-c.ExitBuffChan:
			// conn已关闭
			return
		}
	}
}
// 数据重传
func (c *Connection)Retrans() {
	for{
		if len(c.sendList) <= 0 {
			break
		} else {
			if c.sendList[0].IsBuff {
				c.msgBufChan <- c.sendList[0]
			} else {
				c.msgChan <- c.sendList[0]
			}
			c.sendList = c.sendList[1:]
		}
	}
}

func (c *Connection)SendMsg(data message.IMessage, retrans bool) error{
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data序列化
	propertyList := data.GetPropertyList()
	bts, err := message.WriteByte(propertyList)
	if err != nil {
		fmt.Println("Pack error msg id = ", message.GetMsgId(propertyList))
		return  errors.New("Pack error msg ")
	}

	// 自增消息序号
	// 将消息封包
	c._seq++
	pack := NewDataPack(bts, c._seq, c._ack, false)
	// 添加可靠消息
	if retrans {
		c.sendList = append(c.sendList, pack)
	}

	//写回客户端
	c.msgChan <- pack   //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return nil
}

func (c *Connection)SendBufMsg(data message.IMessage, retrans bool) error{
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}
	//将data封包，并且发送
	propertyList := data.GetPropertyList()
	bts, err := message.WriteByte(propertyList)
	if err != nil {
		fmt.Println("Pack error msg id = ", message.GetMsgId(propertyList))
		return  errors.New("Pack error msg ")
	}

	// 将消息封包
	c._seq++
	pack := NewDataPack(bts, c._seq, c._ack, true)

	// 添加可靠消息
	if retrans {
		c.sendList = append(c.sendList, pack)
	}

	//写回客户端
	c.msgBufChan <- pack   //将之前直接回写给conn.Write的方法 改为 发送给Channel 供Writer读取

	return nil
}

//启动连接，让当前连接开始工作
func (c *Connection) Start() {

	defer c.Stop()

	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()

	hook := c.Server.(ziface.IHook)
	hook.CallOneConnStart(c)

	for {
		select {
		case <- c.ExitBuffChan:
			fmt.Println("得到退出消息，不再阻塞")
			//得到退出消息，不再阻塞
			return
		}
	}
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//fmt.Printf("Stop:%s", debug.Stack())

	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	hook := c.Server.(ziface.IHook)
	hook.CallOneConnStop(c)

	// 关闭socket链接
	c.Conn.Close()

	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true

	c.Server.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgChan)
	close(c.msgBufChan)
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32{
	return c.ConnID
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) GetProperty() ziface.IProperty{
	return c.GetIProperty()
}