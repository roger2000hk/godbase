package godbase

type HashSlots []SkipMap
type HashFn func (Cmp) uint64

type HashMap struct {
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
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *HashMap) Init(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *HashMap {
	m.fn = fn
	m.levels = levels
	m.slots = make(HashSlots, slotCount)

	for _, s := range m.slots {
		s.Init(alloc, levels)
	}

	return m
}

func (m *HashMap) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *HashMap) Len() int64 {
	return m.len
}
