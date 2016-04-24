package maps

type ESkipHashSlots []ESkip

type ESkipHash struct {
	fn HashFn
	len int64
	slots ESkipHashSlots
}

func NewESkipHash(fn HashFn, slotCount int) *ESkipHash {
	return new(ESkipHash).Init(fn, slotCount)
}

func (m *ESkipHash) Cut(start, end Iter, fn TestFn) Any {
	i := m.fn(start.Key()) % uint64(len(m.slots))
	return m.slots[i].Cut(start, end, fn)
}

func (m *ESkipHash) Delete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	i := m.fn(key) % uint64(len(m.slots))
	res, cnt := m.slots[i].Delete(start, end, key, val)
	m.len -= int64(cnt)
	return res, cnt
}

func (m *ESkipHash) Find(start Iter, key Cmp, val interface{}) (Iter, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Find(start, key, val)
	return res, ok
}

func (m *ESkipHash) Init(fn HashFn, slotCount int) *ESkipHash {
	m.fn = fn
	m.slots = make(ESkipHashSlots, slotCount)

	for i, _ := range m.slots {
		m.slots[i].Init()
	}

	return m
}

func (m *ESkipHash) Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(start, key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *ESkipHash) Len() int64 {
	return m.len
}
