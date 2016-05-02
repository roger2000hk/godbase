package sets

import (
	"github.com/fncodr/godbase"
)

type MapHash struct {
	fn godbase.MapHashFn
	len int64
	slotAlloc SlotAlloc
	slots map[interface{}]godbase.Set
}

func (self MapHash) Delete(offs int, key godbase.Key) (godbase.Set, int) {
	si := self.fn(key)

	if s, i := self.slot(key, si).Delete(offs, key); i != -1 {
		self.slots[si] = s.(Sort)
		self.len--
		return self, i
	}

	return self, -1
}

func (self MapHash) Index(offs int, key godbase.Key) int {
	return self.slot(key, self.fn(key)).Index(offs, key)	
}

func (self *MapHash) Init(sc int, fn godbase.MapHashFn, sa SlotAlloc) *MapHash {
	self.fn = fn
	self.slotAlloc = sa
	self.slots = make(map[interface{}]godbase.Set, sc)
	return self
}

func (self MapHash) Insert(offs int, key godbase.Key) (godbase.Set, int) {
	si := self.fn(key)

	if s, i := self.slot(key, si).Insert(offs, key); i != -1 {
		self.slots[si] = s
		self.len++
		return self, i
	}
	
	return self, -1
}

func (self MapHash) Len() int64 {
	return self.len
}

func (self MapHash) slot(key godbase.Key, i interface{}) godbase.Set {
	s, ok := self.slots[i]

	if !ok {
		s = self.slotAlloc(key)
		self.slots[i] = s
	}

	return s
}
