package sets

import (
	"github.com/fncodr/godbase"
)

type Map map[interface{}]bool

func NewMap(s int) Map {
	return make(Map, s)
}

func (self Map) Clone() godbase.Set {
	panic("Clone() not supported")
}

func (self Map) Delete(_ int, k godbase.Key) int {
	if _, ok := self[k]; ok {
		delete(self, k)
		return 1
	}

	return -1
}

func (self Map) Index(_ int, k godbase.Key) int {
	if _, ok := self[k]; ok {
		return 1
	}

	return -1
}

func (self Map) Insert(_ int, k godbase.Key, multi bool) (int, bool) {
	if multi {
		panic("multi not supported")
	}

	_, ok := self[k]

	if ok {
		return 1, false
	}

	self[k] = true
	return 1, true
}

func (self Map) Len() int64 {
	return int64(len(self))
}
