package sets

import (
	"github.com/fncodr/godbase"
)

type Hash struct {
	fn godbase.HashFn
	len int64
	slots []Sort
}

func (self Hash) Delete(offs int, key godbase.Key) (godbase.Set, bool) {
	i := self.fn(key) % uint64(len(self.slots))

	if s, ok := self.slots[i].Delete(offs, key); ok {
		self.slots[i] = s.(Sort)
		self.len--
		return self, true
	}

	return self, false
}

func (self Hash) HasKey(offs int, key godbase.Key) bool {
	i := self.fn(key) % uint64(len(self.slots))
	return self.slots[i].HasKey(offs, key)
}

func (self Hash) Index(offs int, key godbase.Key) int {
	i := self.fn(key) % uint64(len(self.slots))
	return self.slots[i].Index(offs, key)	
}

func (self *Hash) Init(sc int, fn godbase.HashFn) *Hash {
	self.fn = fn
	self.slots = make([]Sort, sc)
	return self
}

func (self Hash) Insert(offs int, key godbase.Key) (godbase.Set, bool) {
	i := self.fn(key) % uint64(len(self.slots))

	if s, ok := self.slots[i].Insert(offs, key); ok {
		self.slots[i] = s.(Sort)
		self.len++
		return self, true
	}
	
	return self, false
}

func (self Hash) Len() int64 {
	return self.len
}
