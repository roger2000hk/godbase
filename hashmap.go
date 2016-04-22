package godbase

type HashSlots []*SkipMap
type HashFn func (Cmp) uint64

type HashMap struct {
	alloc *SkipNodeAlloc
	fn HashFn
	len int64
	levels int
	slots HashSlots
}

func NewHashMap(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *HashMap {
	return new(HashMap).Init(fn, slotCount, alloc, levels)
}

func (m *HashMap) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))

	if s := m.slots[i]; s != nil {
		res := s.Delete(key, val)
		m.len -= int64(res)
		return res
	}

	return 0
}

func (m *HashMap) Init(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *HashMap {
	m.alloc = alloc
	m.fn = fn
	m.levels = levels
	m.slots = make(HashSlots, slotCount)
	return m
}

func (m *HashMap) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.getSlot(key).Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *HashMap) Len() int64 {
	return m.len
}

func (m *HashMap) getSlot(key Cmp) *SkipMap {
	i := m.fn(key) % uint64(len(m.slots))
	s := m.slots[i]

	if s == nil {
		s = NewSkipMap(m.alloc, m.levels)
		m.slots[i] = s
	}

	return s
}
