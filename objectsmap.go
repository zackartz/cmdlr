package cmdlr

import (
	"sync"
)

type ObjectsMap struct {
	mutex    sync.RWMutex
	innerMap map[string]interface{}
}

func NewObjectsMap() *ObjectsMap {
	return &ObjectsMap{
		innerMap: map[string]interface{}{},
	}
}

func (om *ObjectsMap) Get(key string) (interface{}, bool) {
	om.mutex.RLock()
	defer om.mutex.RUnlock()

	v, ok := om.innerMap[key]
	return v, ok
}

func (om *ObjectsMap) MustGet(key string) interface{} {
	v, ok := om.Get(key)
	if !ok {
		return nil
	}
	return v
}

func (om *ObjectsMap) Set(key string, val interface{}) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	om.innerMap[key] = val
}

func (om *ObjectsMap) Delete(key string) {
	om.mutex.Lock()
	defer om.mutex.Unlock()

	delete(om.innerMap, key)
}
