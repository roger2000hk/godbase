package maps

type HashSlots []Skip
type HashFn func (Cmp) uint64

type Hash struct {
	alloc *SkipNodeAlloc
	fn HashFn
	len int64
	levels int
	slots HashSlots
}

func NewHash(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *Hash {
	return new(Hash).Init(fn, slotCount, alloc, levels)
}

func (m *Hash) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *Hash) Init(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *Hash {
	m.alloc = alloc
	m.fn = fn
	m.levels = levels
	m.slots = make(HashSlots, slotCount)
	for i, _ := range m.slots {
		m.slots[i].Init(alloc, levels)
	}
	return m
}

func (m *Hash) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *Hash) Len() int64 {
	return m.len
}
