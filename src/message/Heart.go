package message
import "utils"
/*
	作者 ***
	日期：2019/10/17 20:08:40
	功能：心跳消息
*/
type Heart struct {
	MsgId uint32	// 客户端消息的协议id
	ServerProto byte	// 主协议:ServerProto.go定义的协议
	Req byte	// 心跳识别码
}

func NewHeart(req byte) *Heart {
	var res = new(Heart)
	res.MsgId = MSG_HEART
	res.ServerProto = GameClient_Proto
	res.Req = req
	return res
}

func (this *Heart)GetPropertyList() []interface{} {
	var PropertyList = make([]interface{}, 0)
	PropertyList = append(PropertyList, this.MsgId)
	PropertyList = append(PropertyList, this.ServerProto)
	PropertyList = append(PropertyList, this.Req)
	return PropertyList
}

func (this *Heart)WritePropertyList(PropertyList []interface{}) {
	this.MsgId = uint32(PropertyList[1].(int32))
	this.ServerProto = PropertyList[2].(byte)
	this.Req = PropertyList[3].(byte)
}

func (this *Heart)ToLuaString() string{
	var res = "{"
	res += "MsgId=" + utils.ToString(this.MsgId) + ","
	res += "ServerProto=" + utils.ToString(this.ServerProto) + ","
	res += "Req=" + utils.ToString(this.Req) + ","
	res += "}"
	return res
}