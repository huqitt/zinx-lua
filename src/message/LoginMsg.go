package message
import "utils"
/*
	作者 ***
	日期：2019/10/17 20:04:23
	功能：登录返回消息
*/
type LoginMsg struct {
	MsgId uint32	// 客户端消息的协议id
	ServerProto byte	// 服务器Id
	LoginType byte	// 微信登录为1 QQ登录为2
	Id string	// 登录ID
}

func NewLoginMsg(loginType byte, id string) *LoginMsg {
	var res = new(LoginMsg)
	res.MsgId = MSG_LOGIN
	res.ServerProto = GameClient_Proto
	res.LoginType = loginType
	res.Id = id
	return res
}

func (this *LoginMsg)GetPropertyList() []interface{} {
	var PropertyList = make([]interface{}, 0)
	PropertyList = append(PropertyList, this.MsgId)
	PropertyList = append(PropertyList, this.ServerProto)
	PropertyList = append(PropertyList, this.LoginType)
	PropertyList = append(PropertyList, this.Id)
	return PropertyList
}

func (this *LoginMsg)WritePropertyList(PropertyList []interface{}) {
	this.MsgId = uint32(PropertyList[1].(int32))
	this.ServerProto = PropertyList[2].(byte)
	this.LoginType = PropertyList[3].(byte)
	this.Id = PropertyList[4].(string)
}

func (this *LoginMsg)ToLuaString() string{
	var res = "{"
	res += "MsgId=" + utils.ToString(this.MsgId) + ","
	res += "ServerProto=" + utils.ToString(this.ServerProto) + ","
	res += "LoginType=" + utils.ToString(this.LoginType) + ","
	res += "Id='" + utils.ToString(this.Id) + "',"
	res += "}"
	return res
}