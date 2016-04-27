package maps

import (
	"fmt"
)

type Slots interface {	
	Get(key Key, create bool) Any
	New() Slots
}

type AnySlots struct {
	fn HashFn
	alloc SlotAlloc
	slots []Any
}

type ESkipSlots struct {
	fn HashFn
	slots []ESkip
}

type MapSlots struct {
	fn MapHashFn
	alloc SlotAlloc
	slots map[interface{}]Any
}

type Hash struct {
	isInit bool
	len int64
	slots Slots
}

type HashSlots struct {
	fn HashFn
	alloc SlotsAlloc
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

func NewESkipSlots(sc int, fn HashFn) *ESkipSlots {
	ss := new(ESkipSlots)
	ss.fn = fn
	ss.slots = make([]ESkip, sc)
	return ss
}

func NewHash(slots Slots) *Hash {
	return new(Hash).Init(slots)
}

func NewHashSlots(sc int, fn HashFn, a SlotsAlloc) *HashSlots {
	ss := new(HashSlots)
	ss.fn = fn
	ss.alloc = a
	ss.slots = make([]Hash, sc)
	return ss
}

func NewMapSlots(sc int, fn MapHashFn, a SlotAlloc) *MapSlots {
	ss := new(MapSlots)
	ss.fn = fn
	ss.alloc = a
	ss.slots = make(map[interface{}]Any, sc)
	return ss
}

func NewSkipSlots(sc int, fn HashFn, a *SkipAlloc, ls int) *SkipSlots {
	ss := new(SkipSlots)
	ss.alloc = a
	ss.fn = fn
	ss.levels = ls
	ss.slots = make([]Skip, sc)
	return ss
}

func NewSlots(sc int, fn HashFn, a SlotAlloc) *AnySlots {
	ss := new(AnySlots)
	ss.fn = fn
	ss.alloc = a
	ss.slots = make([]Any, sc)
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

func (m *Hash) First() Iter {
	panic("Hash doesn't support First()!")
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
		s = ss.alloc(key)
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
		return s.Init(ss.alloc(key))
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
		s = ss.alloc(key)
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

func (ss *AnySlots) New() Slots {
	return NewSlots(len(ss.slots), ss.fn, ss.alloc)
}

func (m *Hash) New() Any {
	return NewHash(m.slots.New())
}

func (ss *HashSlots) New() Slots {
	return NewHashSlots(len(ss.slots), ss.fn, ss.alloc)
}

func (ss *MapSlots) New() Slots {
	return NewMapSlots(len(ss.slots), ss.fn, ss.alloc)
}

func (ss *SkipSlots) New() Slots {
	return NewSkipSlots(len(ss.slots), ss.fn, ss.alloc, ss.levels)
}

func (m *Hash) Set(key Key, val interface{}) bool {
	if m.slots.Get(key, true).Set(key, val) {
		m.len++
		return true
	}

	return false
}

func (m *Hash) String() string {
	return fmt.Sprintf("%v", m.slots)
}
