package maps

import (
	"fmt"
)

type Slots interface {	
	Get(key Key, create bool) Any
}

type AnySlots struct {
	fn HashFn
	slotAlloc SlotAlloc
	slots []Any
}

type ESkipSlots struct {
	fn HashFn
	slots []ESkip
}

type MapSlots struct {
	fn MapHashFn
	slotAlloc SlotAlloc
	slots map[interface{}]Any
}

type Hash struct {
	isInit bool
	len int64
	slots Slots
}

type HashSlots struct {
	fn HashFn
	slotsAlloc SlotsAlloc
	slots []Hash
}

type SkipSlots struct {
	alloc *SkipAlloc
	fn HashFn
	levels int
	slots []Skip
}

type HashFn func (Key) uint64
type MapHashFn func (Key) interface{}
type SlotAlloc func (key Key) Any
type SlotsAlloc func (key Key) Slots

func NewESkipSlots(count int, fn HashFn) *ESkipSlots {
	ss := new(ESkipSlots)
	ss.fn = fn
	ss.slots = make([]ESkip, count)
	return ss
}

func NewHash(slots Slots) *Hash {
	return new(Hash).Init(slots)
}

func NewHashSlots(count int, fn HashFn, slotsAlloc SlotsAlloc) *HashSlots {
	ss := new(HashSlots)
	ss.fn = fn
	ss.slotsAlloc = slotsAlloc
	ss.slots = make([]Hash, count)
	return ss
}

func NewMapSlots(count int, fn MapHashFn, slotAlloc SlotAlloc) *MapSlots {
	ss := new(MapSlots)
	ss.fn = fn
	ss.slotAlloc = slotAlloc
	ss.slots = make(map[interface{}]Any, count)
	return ss
}

func NewSkipSlots(count int, fn HashFn, alloc *SkipAlloc, levels int) *SkipSlots {
	ss := new(SkipSlots)
	ss.alloc = alloc
	ss.fn = fn
	ss.levels = levels
	ss.slots = make([]Skip, count)
	return ss
}

func NewSlots(count int, fn HashFn, slotAlloc SlotAlloc) *AnySlots {
	ss := new(AnySlots)
	ss.fn = fn
	ss.slotAlloc = slotAlloc
	ss.slots = make([]Any, count)
	return ss
}

func (m *Hash) Cut(start, end Iter, fn MapFn) Any {
	return m.slots.Get(start.Key(), true).Cut(start, end, fn)
}

func (m *Hash) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	res, cnt := m.slots.Get(key, true).Delete(start, end, key, val)
	m.len -= int64(cnt)
	return res, cnt
}

func (m *Hash) Find(start Iter, key Key, val interface{}) (Iter, bool) {
	return m.slots.Get(key, true).Find(start, key, val)
}

func (m *Hash) Get(key Key) (interface{}, bool) {
	return m.slots.Get(key, true).Get(key)
}

func (ss *AnySlots) Get(key Key, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := ss.slots[i]

	if s != nil {
		return s
	}

	if create {
		s = ss.slotAlloc(key)
		ss.slots[i] = s
		return s
	}

	return nil
}

func (ss *ESkipSlots) Get(key Key, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := &ss.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init()
	}

	return nil
}

func (ss *HashSlots) Get(key Key, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := &ss.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init(ss.slotsAlloc(key))
	}

	return nil
}

func (ss *MapSlots) Get(key Key, create bool) Any {
	i := ss.fn(key)
	s, ok := ss.slots[i]

	if ok {
		return s
	}

	if create {
		s = ss.slotAlloc(key)
		ss.slots[i] = s
		return s
	}

	return nil
}

func (ss *SkipSlots) Get(key Key, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := &ss.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init(ss.alloc, ss.levels)
	}

	return nil
}

func (m *Hash) Init(slots Slots) *Hash {
	m.isInit = true
	m.slots = slots
	return m
}

func (m *Hash) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {	
	res, ok := m.slots.Get(key, true).Insert(start, key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *Hash) Len() int64 {
	return m.len
}

func (m *Hash) String() string {
	return fmt.Sprintf("%v", m.slots)
}
