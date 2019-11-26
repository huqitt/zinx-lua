package ziface

type IProperty interface {
	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
	//
	GetIProperty() IProperty
}
