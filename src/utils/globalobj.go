package utils

import (
	"io/ioutil"
	"encoding/json"
)

type ServerConfig struct {
	Id		   uint32
	Host       string         //当前服务器主机IP
	Name       string         //当前服务器名称
	ServerName string         //当前服务器名称

	/*
		Zinx
	*/
	Version          string //当前Zinx版本号
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen uint32 	//缓冲的发送Chan的存储数量

	/*
		config file path
	*/
	ConfFilePath string

	// RPC
	RegRPC			bool 	// 需要注册RPC监听
	ConnRPC			bool	// 需要链接到RPC服务器
}


type OtherServerConfig struct {
	Id		  		uint32
	Name	  		string			//名字
	Host      		string         	//链接的服务器主机IP、port
	Origin	  		string			//
	ServerName	  	string			//
}

var GlobalServer *ServerConfig
var OtherServerConfigList []OtherServerConfig

//读取用户的配置文件
func (g *ServerConfig) ReloadServerConfig() {
	data, err := ioutil.ReadFile("conf/serverConfig.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &GlobalServer)
	if err != nil {
		panic(err)
	}
}
//读取用户的配置文件
func (g *ServerConfig) ReloadClientConfig() {
	data, err := ioutil.ReadFile("conf/otherServerConfig.json")
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	//fmt.Printf("json :%s\n", data)
	err = json.Unmarshal(data, &OtherServerConfigList)
	if err != nil {
		panic(err)
	}
}

/*
    提供init方法，默认加载
*/
func init() {
	//初始化GlobalObject变量，设置一些默认值
	GlobalServer = &ServerConfig{
		Name:          "ZinxServerApp",
		Version:       "V0.4",
		Host:          "127.0.0.1",
		MaxConn:       12000,
		MaxPacketSize: 4096,
		ConfFilePath:  "conf/serverConfig.json",
		WorkerPoolSize: 10,
		MaxWorkerTaskLen: 1024,
	}

	//从配置文件中加载一些用户配置的参数
	GlobalServer.ReloadServerConfig()
	GlobalServer.ReloadClientConfig()
}
