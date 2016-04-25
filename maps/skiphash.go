package maps

type SkipHashSlots []Skip
type HashFn func (Cmp) uint64

type SkipHash struct {
	alloc *SkipAlloc
	fn HashFn
	len int64
	levels int
	slots SkipHashSlots
}

func NewSkipHash(fn HashFn, slotCount int, alloc *SkipAlloc, levels int) *SkipHash {
	return new(SkipHash).Init(fn, slotCount, alloc, levels)
}

func (m *SkipHash) Cut(start, end Iter, fn TestFn) Any {
	i := m.fn(start.Key()) % uint64(len(m.slots))
	return m.slots[i].Cut(start, end, fn)
}

func (m *SkipHash) Delete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	i := m.fn(key) % uint64(len(m.slots))
	res, cnt := m.slots[i].Delete(start, end, key, val)
	m.len -= int64(cnt)
	return res, cnt
}

func (m *SkipHash) Find(start Iter, key Cmp, val interface{}) (Iter, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Find(start, key, val)
	return res, ok
}

func (m *SkipHash) Init(fn HashFn, slotCount int, alloc *SkipAlloc, levels int) *SkipHash {
	m.alloc = alloc
	m.fn = fn
	m.levels = levels
	m.slots = make(SkipHashSlots, slotCount)
	for i, _ := range m.slots {
		m.slots[i].Init(alloc, levels)
	}
	return m
}

func (m *SkipHash) Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {	
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(start, key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *SkipHash) Len() int64 {
	return m.len
}
