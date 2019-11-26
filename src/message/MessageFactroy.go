package message

const(
	GameClient_Proto	byte = 2	//	客户端的协议
	GameDataDB_Proto	byte = 3	//	游戏的DB主协议
	GameNet_Proto	byte = 4	//	游戏的Net主协议
	G_Error_Proto	byte = 5	//	游戏的错误处理
	G_GateWay_Proto	byte = 6	//	网关协议
	G_GameHall_Proto	byte = 7	//	大厅协议
	G_GameLogin_Proto	byte = 8	//	登录服务器协议
	G_GameGlobal_Proto	byte = 9	//	负责全局的游戏逻辑
)

const(
	MSG_HEART	uint32 = 1	//	心跳消息
	MSG_LUA_STRING	uint32 = 2	//	lua字符串消息
	MSG_SER_ID	uint32 = 3	//	服务器ID
	MSG_LOGIN	uint32 = 4	//	
)

func GetMessage(propertyList []interface{}) IMessage {
	var proto = uint32(propertyList[1].(int32))
	switch proto {
	case MSG_LUA_STRING:
		var message = &LuaMessage{}
		message.WritePropertyList(propertyList)
		return message
	}
	return nil
}
