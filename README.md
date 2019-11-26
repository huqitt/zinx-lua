# zinx-lua
zinx-lua

## 说明

这是一个通信采用GO语言编写，服务逻辑采用lua编写的服务器框架，是zinx的服务器的分支。

## 安装

### 1.克隆下载源码
>>>https://github.com/huqitt/zinx-lua

### 2.下载Go-lua插件（此步骤可以省略）
>>>https://github.com/yuin/gopher-lua

### 3.推荐使用IDEA打开




## 启动

### 1.配置conf下serverConfig.json文件（可以自定义lua虚拟机最大数量）

### 2.src下Script下main.lua配置提供消息接口：

>>（1）OnMessage（ser_id, table, func, param）收到服务器消息时触发

>>>ser_id:发消息的服务器id
  
>>>table：字符串，要调用的表的名字
  
>>>func：字符串，要调用表的方法
  
>>>param：是一个table，用于传递参数
  
>>（2）Send（serverId, table, func, param）发送消息到其他服务器

>>>serverId：接受消息的服务器id；
  
>>>table：字符串，调用目标服务器的table的名字
  
>>>func：字符串类型，调用的目标服务器的table定义的方法
  
>>>param：table类型，传给目标服务器的参数
  
>>（3）AddServer（ser_id, name, desc）有服务器连接到本机时触发
>>（4）编写业务逻辑或者lua脚本处理服务器消息
