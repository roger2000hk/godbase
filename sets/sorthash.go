package sets

import (
	"github.com/fncodr/godbase"
)

type SortHash struct {
	fn godbase.HashFn
	len int64
	slots []Sort
}

func (self SortHash) Delete(offs int, key godbase.Key) (godbase.Set, int) {
	si := self.fn(key) % uint64(len(self.slots))

	if s, i := self.slots[si].Delete(offs, key); i != -1 {
		self.slots[si] = s.(Sort)
		self.len--
		return self, i
	}

	return self, -1
}

func (self SortHash) Index(offs int, key godbase.Key) int {
	si := self.fn(key) % uint64(len(self.slots))
	return self.slots[si].Index(offs, key)	
}

func (self *SortHash) Init(sc int, fn godbase.HashFn) *SortHash {
	self.fn = fn
	self.slots = make([]Sort, sc)
	return self
}

func (self SortHash) Insert(offs int, key godbase.Key) (godbase.Set, int) {
	si := self.fn(key) % uint64(len(self.slots))

	if s, i := self.slots[si].Insert(offs, key); i != -1 {
		self.slots[si] = s.(Sort)
		self.len++
		return self, i
	}
	
	return self, -1
}

func (self SortHash) Len() int64 {
	return self.len
}
