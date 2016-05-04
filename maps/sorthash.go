package maps

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type SortHash struct {
	alloc *SlabAlloc
	fn godbase.HashFn
	isInit bool
	len int64
	levels int
	slots []Sort
}

func NewSlabHash(sc int, fn godbase.HashFn, a *SlabAlloc, ls int) *SortHash {
	return new(SortHash).Init(sc, fn, a, ls)
}

func NewSortHash(sc int, fn godbase.HashFn, ls int) *SortHash {
	return NewSlabHash(sc, fn, nil, ls)
}

func (self *SortHash) Init(sc int, fn godbase.HashFn, a *SlabAlloc, ls int) *SortHash {
	self.alloc = a
	self.fn = fn
	self.levels = ls
	self.slots = make([]Sort, sc)
	return self
}

func (self *SortHash) Clear() {
	for i := range self.slots {
		self.slots[i].Clear()
	}
	
	self.len = 0
}

func (self *SortHash) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	return self.GetSlot(start.Key(), true).Cut(start, end, fn)
}

func (self *SortHash) Delete(start, end godbase.Iter, key godbase.Key, 
	val interface{}) (godbase.Iter, int) {
	res, cnt := self.GetSlot(key, true).Delete(start, end, key, val)
	self.len -= int64(cnt)
	return res, cnt
}

func (self *SortHash) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return self.GetSlot(key, true).Find(start, key, val)
}

func (_ *SortHash) First() godbase.Iter {
	panic("First() not supported")
}

func (self *SortHash) Get(key godbase.Key) (interface{}, bool) {
	return self.GetSlot(key, true).Get(key)
}

func (self *SortHash) GetSlot(key godbase.Key, create bool) godbase.Map {
	i := self.fn(key) % uint64(len(self.slots))
	s := &self.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init(self.alloc, self.levels)
	}

	return nil
}

func (self *SortHash) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	res, ok := self.GetSlot(key, true).Insert(start, key, val, allowMulti)

	if ok {
		self.len++
	}

	return res, ok
}

func (self *SortHash) Len() int64 {
	return self.len
}

func (self *SortHash) New() godbase.Map {
	return NewSlabHash(len(self.slots), self.fn, self.alloc, self.levels)
}

func (self *SortHash) Set(key godbase.Key, val interface{}) bool {
	if self.GetSlot(key, true).Set(key, val) {
		self.len++
		return true
	}

	return false
}

func (self *SortHash) String() string {
	return fmt.Sprintf("%v", self.slots)
}

func (self *SortHash) While(fn godbase.KVTestFn) bool {
	for _, s := range self.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true	
}
