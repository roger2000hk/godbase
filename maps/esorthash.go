package maps

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type ESortHash struct {
	fn godbase.HashFn
	isInit bool
	len int64
	slots []ESort
}

func NewESortHash(sc int, fn godbase.HashFn) *ESortHash {
	return new(ESortHash).Init(sc, fn)
}

func (self *ESortHash) Clear() {
	for i := range self.slots {
		self.slots[i].Clear()
	}
	
	self.len = 0
}

func (self *ESortHash) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	return self.GetSlot(start.Key(), true).Cut(start, end, fn)
}

func (self *ESortHash) Delete(start, end godbase.Iter, key godbase.Key, 
	val interface{}) (godbase.Iter, int) {
	res, cnt := self.GetSlot(key, true).Delete(start, end, key, val)
	self.len -= int64(cnt)
	return res, cnt
}

func (self *ESortHash) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return self.GetSlot(key, true).Find(start, key, val)
}

func (_ *ESortHash) First() godbase.Iter {
	panic("First() not supported")
}

func (self *ESortHash) Get(key godbase.Key) (interface{}, bool) {
	return self.GetSlot(key, true).Get(key)
}

func (self *ESortHash) GetSlot(key godbase.Key, create bool) godbase.Map {
	i := self.fn(key) % uint64(len(self.slots))
	s := &self.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init()
	}

	return nil
}

func (self *ESortHash) Init(sc int, fn godbase.HashFn) *ESortHash {
	self.fn = fn
	self.slots = make([]ESort, sc)
	self.isInit = true
	return self
}

func (self *ESortHash) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	res, ok := self.GetSlot(key, true).Insert(start, key, val, allowMulti)

	if ok {
		self.len++
	}

	return res, ok
}

func (self *ESortHash) Len() int64 {
	return self.len
}

func (self *ESortHash) New() godbase.Map {
	return NewESortHash(len(self.slots), self.fn)
}

func (self *ESortHash) Set(key godbase.Key, val interface{}) bool {
	if self.GetSlot(key, true).Set(key, val) {
		self.len++
		return true
	}

	return false
}

func (self *ESortHash) String() string {
	return fmt.Sprintf("%v", self.slots)
}

func (self *ESortHash) While(fn godbase.KVTestFn) bool {
	for _, s := range self.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true	
}
