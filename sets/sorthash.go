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

func (self *SortHash) Delete(offs int, key godbase.Key) int {
	si := self.fn(key) % uint64(len(self.slots))

	if i := self.slots[si].Delete(offs, key); i != -1 {
		self.len--
		return i
	}

	return -1
}

func (self *SortHash) Index(offs int, key godbase.Key) int {
	si := self.fn(key) % uint64(len(self.slots))
	return self.slots[si].Index(offs, key)	
}

func (self *SortHash) Init(sc int, fn godbase.HashFn) *SortHash {
	self.fn = fn
	self.slots = make([]Sort, sc)
	return self
}

func (self *SortHash) Insert(offs int, key godbase.Key) int {
	si := self.fn(key) % uint64(len(self.slots))

	if i := self.slots[si].Insert(offs, key); i != -1 {
		self.len++
		return i
	}
	
	return -1
}

func (self *SortHash) Len() int64 {
	return self.len
}
