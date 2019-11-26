package message

type IMessage interface {
	GetPropertyList() []interface{}
	WritePropertyList([]interface{})
	ToLuaString() string
}
