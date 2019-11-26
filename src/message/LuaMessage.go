package message
import "utils"
/*
	作者 ***
	日期：2019/10/17 20:09:10
	功能：lua消息
*/
type LuaMessage struct {
	MsgId uint32	// 消息id
	ServerProto byte	// ServerId
	Table string	// 脚本名字
	Function string	// 方法名字
	Param string	// 参数
}

func NewLuaMessage(table string, function string, param string) *LuaMessage {
	var res = new(LuaMessage)
	res.MsgId = MSG_LUA_STRING
	res.ServerProto = GameClient_Proto
	res.Table = table
	res.Function = function
	res.Param = param
	return res
}

func (this *LuaMessage)GetPropertyList() []interface{} {
	var PropertyList = make([]interface{}, 0)
	PropertyList = append(PropertyList, this.MsgId)
	PropertyList = append(PropertyList, this.ServerProto)
	PropertyList = append(PropertyList, this.Table)
	PropertyList = append(PropertyList, this.Function)
	PropertyList = append(PropertyList, this.Param)
	return PropertyList
}

func (this *LuaMessage)WritePropertyList(PropertyList []interface{}) {
	this.MsgId = uint32(PropertyList[1].(int32))
	this.ServerProto = PropertyList[2].(byte)
	this.Table = PropertyList[3].(string)
	this.Function = PropertyList[4].(string)
	this.Param = PropertyList[5].(string)
}

func (this *LuaMessage)ToLuaString() string{
	var res = "{"
	res += "MsgId=" + utils.ToString(this.MsgId) + ","
	res += "ServerProto=" + utils.ToString(this.ServerProto) + ","
	res += "Table='" + utils.ToString(this.Table) + "',"
	res += "Function='" + utils.ToString(this.Function) + "',"
	res += "Param='" + utils.ToString(this.Param) + "',"
	res += "}"
	return res
}