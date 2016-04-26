package maps

import (
	"fmt"
)

type Map map[Key]interface{}

func NewMap() Map {
	return make(Map)
}

func (m Map) Cut(start, end Iter, fn MapFn) Any {
	panic("Map doesn't support iters")
}

func (m Map) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	if start != nil || end != nil {
		panic("Map doesn't support iters")
	}

	if val != nil {
		panic("Map doesn't support multi")
	}

	delete(m, key)
	return nil, 1
}

func (m Map) Find(start Iter, key Key, val interface{}) (Iter, bool) {
	v, ok := m[key]
	return nil, ok && (val == nil || v == val)
}

func (m Map) Get(key Key) (interface{}, bool) {
	v, ok := m[key]
	return v, ok
}

func (m Map) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	if start != nil {
		panic("Map doesn't support iters")
	}

	if allowMulti {
		panic("Map doesn't support multi")
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

func (m Map) Set(key Key, val interface{}) interface{} {
	m[key] = val
	return val
}

func (m Map) String() string {
	return fmt.Sprintf("%v", m)
}
