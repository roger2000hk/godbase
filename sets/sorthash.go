package sets

import (
	"github.com/fncodr/godbase"
)

type SortHash struct {
	fn godbase.HashFn
	len int64
	slots []Sort
}

func (self SortHash) Delete(offs int, key godbase.Key) (godbase.Set, bool) {
	i := self.fn(key) % uint64(len(self.slots))

	if s, ok := self.slots[i].Delete(offs, key); ok {
		self.slots[i] = s.(Sort)
		self.len--
		return self, true
	}

	return self, false
}

func (self SortHash) HasKey(offs int, key godbase.Key) bool {
	i := self.fn(key) % uint64(len(self.slots))
	return self.slots[i].HasKey(offs, key)
}

func (self SortHash) Index(offs int, key godbase.Key) int {
	i := self.fn(key) % uint64(len(self.slots))
	return self.slots[i].Index(offs, key)	
}

func (self *SortHash) Init(sc int, fn godbase.HashFn) *SortHash {
	self.fn = fn
	self.slots = make([]Sort, sc)
	return self
}

func (self SortHash) Insert(offs int, key godbase.Key) (godbase.Set, bool) {
	i := self.fn(key) % uint64(len(self.slots))

	if s, ok := self.slots[i].Insert(offs, key); ok {
		self.slots[i] = s.(Sort)
		self.len++
		return self, true
	}
	
	return self, false
}

func (self SortHash) Len() int64 {
	return self.len
}
