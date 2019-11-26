package message
import "utils"
/*
	作者 ***
	日期：2019/10/17 20:08:40
	功能：lua消息
*/
type ServerIdMsg struct {
	MsgId uint32	// 消息id
	ServerProto byte	// ServerId
	ServerId uint32	// 服务器Id
	Name string	// 名字
	Desc string	// 备注
}

func NewServerIdMsg(serverId uint32, name string, desc string) *ServerIdMsg {
	var res = new(ServerIdMsg)
	res.MsgId = MSG_SER_ID
	res.ServerProto = GameClient_Proto
	res.ServerId = serverId
	res.Name = name
	res.Desc = desc
	return res
}

func (this *ServerIdMsg)GetPropertyList() []interface{} {
	var PropertyList = make([]interface{}, 0)
	PropertyList = append(PropertyList, this.MsgId)
	PropertyList = append(PropertyList, this.ServerProto)
	PropertyList = append(PropertyList, this.ServerId)
	PropertyList = append(PropertyList, this.Name)
	PropertyList = append(PropertyList, this.Desc)
	return PropertyList
}

func (this *ServerIdMsg)WritePropertyList(PropertyList []interface{}) {
	this.MsgId = uint32(PropertyList[1].(int32))
	this.ServerProto = PropertyList[2].(byte)
	this.ServerId = uint32(PropertyList[3].(int32))
	this.Name = PropertyList[4].(string)
	this.Desc = PropertyList[5].(string)
}

func (this *ServerIdMsg)ToLuaString() string{
	var res = "{"
	res += "MsgId=" + utils.ToString(this.MsgId) + ","
	res += "ServerProto=" + utils.ToString(this.ServerProto) + ","
	res += "ServerId=" + utils.ToString(this.ServerId) + ","
	res += "Name='" + utils.ToString(this.Name) + "',"
	res += "Desc='" + utils.ToString(this.Desc) + "',"
	res += "}"
	return res
}