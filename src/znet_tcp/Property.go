package znet_tcp

import (
	"sync"
	"errors"
	"ziface"
)

type Property struct {
	property map[string]interface{}
	lock	 sync.RWMutex
}

func (p *Property)GetIProperty() ziface.IProperty{
	return p
}

func (p *Property) SetProperty(key string, value interface{}){
	p.lock.Lock()
	defer p.lock.Unlock()

	p.property[key] = value
}

func (p *Property) GetProperty(key string)(interface{}, error){
	p.lock.Lock()
	defer p.lock.Unlock()

	value,ok := p.property[key]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("key " + key + " not found!")
	}
}

func (p *Property) RemoveProperty(key string) {
	p.lock.Lock()
	defer p.lock.Unlock()

	delete(p.property, key)
}