# zinx-lua
zinx-lua

1.目录结构

zinx
--conf(服务器配置文件)
--src（代码目录）
  --LuaState（初始化lua虚拟机）
  --main（服务器启动、客户端测试）
  --message(消息结构体)
  --Rotuter(消息处理)
    --LuaMsgRouter
    --ServerIdRouter
  --Script（lua代码文件）
  --utils
  --ziface
  --znet_tcp
