package utils

import (
	"sync"
)

type SafeDict struct {
	Data map[string]interface{}
	*sync.RWMutex
}

func NewSafeDict(data map[string]interface{}) *SafeDict {
	return &SafeDict{data, &sync.RWMutex{}}
}

func (d *SafeDict) Len() int {
	d.RLock()
	defer d.RUnlock()
	return len(d.Data)
}

func (d *SafeDict) Put(key string, value interface{}) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.Data[key]
	d.Data[key] = value
	return oldValue, ok
}

func (d *SafeDict) Get(key string) (interface{}, bool) {
	d.RLock()
	defer d.RUnlock()
	oldValue, ok := d.Data[key]
	return oldValue, ok
}

func (d *SafeDict) Exists(key string) bool {
	d.RLock()
	defer d.RUnlock()
	_, ok := d.Data[key]
	if ok {
		return true
	} else {
		return false
	}
}

func (d *SafeDict) Delete(key string) (interface{}, bool) {
	d.Lock()
	defer d.Unlock()
	oldValue, ok := d.Data[key]
	if ok {
		delete(d.Data, key)
	}
	return oldValue, ok
}

type SyncMap struct {
	sync.Map
}

func NewSyncMap(data map[string]interface{}) *SyncMap {
	var syncMap SyncMap
	for idx, val := range data {
		syncMap.Store(idx, val)
	}
	return &syncMap
}

func (sm *SyncMap) ToMap() map[string]interface{} {
	m := map[string]interface{}{}
	sm.Range(func(key, value interface{}) bool {
		name := key.(string)
		m[name] = value
		return true
	})
	return m
}
