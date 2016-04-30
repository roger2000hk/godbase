package maps

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type Map map[godbase.Key]interface{}

func NewMap() Map {
	return make(Map)
}

func (m Map) Clear() {
	for k := range m {
		delete(m, k)
	}
}

func (m Map) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	panic("Map doesn't support iters")
}

func (m Map) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	if start != nil || end != nil {
		panic("Map doesn't support iters")
	}

	if val != nil {
		panic("Map doesn't support multi")
	}

	delete(m, key)
	return nil, 1
}

func (m Map) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	v, ok := m[key]
	return nil, ok && (val == nil || v == val)
}

func (m Map) First() godbase.Iter {
	panic("Map doesn't support iters!")
}

func (m Map) Get(key godbase.Key) (interface{}, bool) {
	v, ok := m[key]
	return v, ok
}

func (m Map) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	if start != nil {
		panic("Map doesn't support iters!")
	}

	if allowMulti {
		panic("Map doesn't support multi!")
	}

	if _, ok := m[key]; ok {
		return nil, false
	}
	
	m[key] = val
	return nil, true
}

func (m Map) Len() int64 {
	return int64(len(m))
}

func (m Map) New() godbase.Map {
	return NewMap()
}

func (m Map) Set(key godbase.Key, val interface{}) bool {
	_, ok := m[key]
	m[key] = val
	return ok
}

func (m Map) String() string {
	return fmt.Sprintf("%v", m)
}

func (m Map) While(fn godbase.KVTestFn) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}
	
	return true
}
