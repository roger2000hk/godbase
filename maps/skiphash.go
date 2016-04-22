package maps

type SkipHashSlots []Skip
type HashFn func (Cmp) uint64

type SkipHash struct {
	alloc *SkipNodeAlloc
	fn HashFn
	len int64
	levels int
	slots SkipHashSlots
}

func NewSkipHash(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *SkipHash {
	return new(SkipHash).Init(fn, slotCount, alloc, levels)
}

func (m *SkipHash) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *SkipHash) Init(fn HashFn, slotCount int, alloc *SkipNodeAlloc, levels int) *SkipHash {
	m.alloc = alloc
	m.fn = fn
	m.levels = levels
	m.slots = make(SkipHashSlots, slotCount)
	for i, _ := range m.slots {
		m.slots[i].Init(alloc, levels)
	}
	return m
}

func (m *SkipHash) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *SkipHash) Len() int64 {
	return m.len
}
