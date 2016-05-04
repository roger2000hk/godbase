package sets

import (
	"github.com/fncodr/godbase"
)

type SortSlots []Sort

type SortHash struct {
	fn godbase.HashFn
	len int64
	slots SortSlots
}

func NewSortHash(sc int, fn godbase.HashFn) *SortHash {
	return new(SortHash).Init(sc, fn)
}

func (self *SortHash) Clone() godbase.Set {
	res := &SortHash{
		fn: self.fn,
		len: self.len,
		slots: make(SortSlots, len(self.slots)) }
	
	for i, s := range self.slots {
		res.slots[i] = *s.Clone().(*Sort)
	}

	return res
}

func (self *SortHash) Delete(start int, key godbase.Key) int {
	si := self.fn(key) % uint64(len(self.slots))

	if i := self.slots[si].Delete(start, key); i != -1 {
		self.len--
		return i
	}

	return -1
}

func (self *SortHash) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	si := self.fn(key) % uint64(len(self.slots))
	i, ok := self.slots[si].DeleteAll(start, end, key)
	self.len -= ok
	return i, ok
}

func (self *SortHash) First(start int, key godbase.Key) (int, bool) {
	si := self.fn(key) % uint64(len(self.slots))
	return self.slots[si].First(start, key)	
}

func (self *SortHash) Get(key godbase.Key, i int) godbase.Key {
	si := self.fn(key) % uint64(len(self.slots))
	return self.slots[si].Get(key, i)
}

func (self *SortHash) Init(sc int, fn godbase.HashFn) *SortHash {
	self.fn = fn
	self.slots = make([]Sort, sc)
	return self
}

func (self *SortHash) Last(start, end int, key godbase.Key) (int, bool) {
	si := self.fn(key) % uint64(len(self.slots))
	return self.slots[si].Last(start, end, key)	
}

func (self *SortHash) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	si := self.fn(key) % uint64(len(self.slots))
	i, ok := self.slots[si].Insert(start, key, multi)
	
	if ok {
		self.len++
	}
	
	return i, ok
}

func (self *SortHash) Len() int64 {
	return self.len
}

func (self *SortHash) While(fn godbase.IKTestFn) bool {
	for _, s := range self.slots {
		if !s.While(fn) {
			return false
		}
	}

	return true
}
