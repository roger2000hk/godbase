package maps

import (
	"fmt"
)

type Slots interface {
	Clear()
	Get(key Key, create bool) Any
	New() Slots
	While(fn TestFn) bool
}

type AnySlots struct {
	fn HashFn
	alloc SlotAlloc
	slots []Any
}

type ESortSlots struct {
	fn HashFn
	slots []ESort
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

type SortSlots struct {
	alloc *SlabAlloc
	fn HashFn
	levels int
	slots []Sort
}

type HashFn func (Key) uint64
type MapHashFn func (Key) interface{}
type SlotAlloc func (key Key) Any
type SlotsAlloc func (key Key) Slots

func NewESortSlots(sc int, fn HashFn) *ESortSlots {
	ss := new(ESortSlots)
	ss.fn = fn
	ss.slots = make([]ESort, sc)
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

func NewSlabSlots(sc int, fn HashFn, a *SlabAlloc, ls int) *SortSlots {
	ss := new(SortSlots)
	ss.alloc = a
	ss.fn = fn
	ss.levels = ls
	ss.slots = make([]Sort, sc)
	return ss
}

func NewSortSlots(sc int, fn HashFn, ls int) *SortSlots {
	return NewSlabSlots(sc, fn, nil, ls)
}

func NewSlots(sc int, fn HashFn, a SlotAlloc) *AnySlots {
	ss := new(AnySlots)
	ss.fn = fn
	ss.alloc = a
	ss.slots = make([]Any, sc)
	return ss
}

func (ss *AnySlots) Clear() {
	for i := range ss.slots {
		ss.slots[i] = nil
	}
}

func (ss *ESortSlots) Clear() {
	for i := range ss.slots {
		ss.slots[i].Clear()
	}
}

func (m *Hash) Clear() {
	m.slots.Clear()
	m.len = 0
}

func (ss *HashSlots) Clear() {
	for i := range ss.slots {
		ss.slots[i].Clear()
	}
}

func (ss *MapSlots) Clear() {
	for i := range ss.slots {
		ss.slots[i] = nil
	}
}

func (ss *SortSlots) Clear() {
	for i := range ss.slots {
		ss.slots[i].Clear()
	}
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

func (ss *ESortSlots) Get(key Key, create bool) Any {
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

func (ss *SortSlots) Get(key Key, create bool) Any {
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

func (ss *ESortSlots) New() Slots {
	return NewESortSlots(len(ss.slots), ss.fn)
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

func (ss *SortSlots) New() Slots {
	return NewSlabSlots(len(ss.slots), ss.fn, ss.alloc, ss.levels)
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

func (ss *AnySlots) While(fn TestFn) bool {
	for _, s := range ss.slots {
		if s != nil && !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *ESortSlots) While(fn TestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

func (m *Hash) While(fn TestFn) bool {
	return m.slots.While(fn)
}

func (ss *HashSlots) While(fn TestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *MapSlots) While(fn TestFn) bool {
	for _, s := range ss.slots {
		if !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *SortSlots) While(fn TestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

