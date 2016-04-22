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

func (m *ESkipHash) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *ESkipHash) Init(fn HashFn, slotCount int) *ESkipHash {
	m.fn = fn
	m.slots = make(ESkipHashSlots, slotCount)

	for i, _ := range m.slots {
		m.slots[i].Init()
	}

	return m
}

func (m *ESkipHash) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *ESkipHash) Len() int64 {
	return m.len
}
