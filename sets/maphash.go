package sets

import (
	"github.com/fncodr/godbase"
)

type MapSlots map[interface{}]godbase.Set

type MapHash struct {
	fn godbase.MapHashFn
	len int64
	slotAlloc SlotAlloc
	slots MapSlots
}

func (self *MapHash) Clone() godbase.Set {
	res := &MapHash{
		fn: self.fn,
		len: self.len,
		slotAlloc: self.slotAlloc,
		slots: make(MapSlots, len(self.slots)) }
	
	for i, s := range self.slots {
		res.slots[i] = s.Clone()
	}

	return res
}

func (self *MapHash) Delete(start int, key godbase.Key) int {
	si := self.fn(key)

	if i := self.slot(key, si).Delete(start, key); i != -1 {
		self.len--
		return i
	}

	return -1
}

func (self *MapHash) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	si := self.fn(key)
	i, ok := self.slot(key, si).DeleteAll(start, end, key)
	self.len -= ok
	return i, ok
}

func (self *MapHash) First(start int, key godbase.Key) int {
	return self.slot(key, self.fn(key)).First(start, key)	
}

func (self *MapHash) Last(start, end int, key godbase.Key) int {
	return self.slot(key, self.fn(key)).Last(start, end, key)	
}

func (self *MapHash) Init(sc int, fn godbase.MapHashFn, sa SlotAlloc) *MapHash {
	self.fn = fn
	self.slotAlloc = sa
	self.slots = make(map[interface{}]godbase.Set, sc)
	return self
}

func (self *MapHash) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	i, ok := self.slot(key, self.fn(key)).Insert(start, key, multi)

	if ok {
		self.len++
	}
	
	return i, ok
}

func (self *MapHash) Len() int64 {
	return self.len
}

func (self *MapHash) slot(key godbase.Key, i interface{}) godbase.Set {
	s, ok := self.slots[i]

	if !ok {
		s = self.slotAlloc(key)
		self.slots[i] = s
	}

	return s
}
