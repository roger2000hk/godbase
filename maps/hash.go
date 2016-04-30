package maps

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type Slots interface {
	Clear()
	Get(key godbase.Key, create bool) godbase.Map
	New() Slots
	While(fn godbase.KVTestFn) bool
}

type BasicSlots struct {
	fn HashFn
	alloc SlotAlloc
	slots []godbase.Map
}

type ESortSlots struct {
	fn HashFn
	slots []ESort
}

type MapSlots struct {
	fn MapHashFn
	alloc SlotAlloc
	slots map[interface{}]godbase.Map
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

type HashFn func (godbase.Key) uint64
type MapHashFn func (godbase.Key) interface{}
type SlotAlloc func (key godbase.Key) godbase.Map
type SlotsAlloc func (key godbase.Key) Slots

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
	ss.slots = make(map[interface{}]godbase.Map, sc)
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

func NewSlots(sc int, fn HashFn, a SlotAlloc) *BasicSlots {
	ss := new(BasicSlots)
	ss.fn = fn
	ss.alloc = a
	ss.slots = make([]godbase.Map, sc)
	return ss
}

func (ss *BasicSlots) Clear() {
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

func (m *Hash) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	return m.slots.Get(start.Key(), true).Cut(start, end, fn)
}

func (m *Hash) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	res, cnt := m.slots.Get(key, true).Delete(start, end, key, val)
	m.len -= int64(cnt)
	return res, cnt
}

func (m *Hash) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return m.slots.Get(key, true).Find(start, key, val)
}

func (m *Hash) First() godbase.Iter {
	panic("Hash doesn't support First()!")
}

func (m *Hash) Get(key godbase.Key) (interface{}, bool) {
	return m.slots.Get(key, true).Get(key)
}

func (ss *BasicSlots) Get(key godbase.Key, create bool) godbase.Map {
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

func (ss *ESortSlots) Get(key godbase.Key, create bool) godbase.Map {
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

func (ss *HashSlots) Get(key godbase.Key, create bool) godbase.Map {
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

func (ss *MapSlots) Get(key godbase.Key, create bool) godbase.Map {
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

func (ss *SortSlots) Get(key godbase.Key, create bool) godbase.Map {
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

func (m *Hash) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	res, ok := m.slots.Get(key, true).Insert(start, key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *Hash) Len() int64 {
	return m.len
}

func (ss *BasicSlots) New() Slots {
	return NewSlots(len(ss.slots), ss.fn, ss.alloc)
}

func (ss *ESortSlots) New() Slots {
	return NewESortSlots(len(ss.slots), ss.fn)
}

func (m *Hash) New() godbase.Map {
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

func (m *Hash) Set(key godbase.Key, val interface{}) bool {
	if m.slots.Get(key, true).Set(key, val) {
		m.len++
		return true
	}

	return false
}

func (m *Hash) String() string {
	return fmt.Sprintf("%v", m.slots)
}

func (ss *BasicSlots) While(fn godbase.KVTestFn) bool {
	for _, s := range ss.slots {
		if s != nil && !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *ESortSlots) While(fn godbase.KVTestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

func (m *Hash) While(fn godbase.KVTestFn) bool {
	return m.slots.While(fn)
}

func (ss *HashSlots) While(fn godbase.KVTestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *MapSlots) While(fn godbase.KVTestFn) bool {
	for _, s := range ss.slots {
		if !s.While(fn) {
			return false
		}
	}

	return true
}

func (ss *SortSlots) While(fn godbase.KVTestFn) bool {
	for _, s := range ss.slots {
		if s.isInit && !s.While(fn) {
			return false
		}
	}

	return true
}

