package sets

import (
	"github.com/fncodr/godbase"
)

type Map map[godbase.Key]bool

func NewMap(s int) Map {
	return make(Map, s)
}

func (self Map) Clone() godbase.Set {
	res := make(Map, len(self))

	for k, _ := range self {
		res[k] = true
	}

	return res
}

func (self Map) Delete(_ int, k godbase.Key) int {
	if _, ok := self[k]; ok {
		delete(self, k)
		return 1
	}

	return -1
}

func (self Map) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	panic("DeleteAll() not supported")
}

func (self Map) First(_ int, k godbase.Key) (int, bool) {
	if _, ok := self[k]; ok {
		return 1, true
	}

	return -1, false
}

func (self Map) Get(k godbase.Key, i int) godbase.Key {
	panic("Get() not supported")
}

func (self Map) Last(start, _ int, k godbase.Key) (int, bool) {
	return self.First(start, k)
}

func (self Map) Load(_ int, keys...godbase.Key) {
	for _, k := range keys {
		self[k] = true
	}
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

func (self Map) While(fn godbase.IKTestFn) bool {
	for k, _ := range self {
		if !fn(-1, k) {
			return false
		}
	}
	
	return true
}
