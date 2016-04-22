package maps

type EHashSlots []ESkip

type EHash struct {
	fn HashFn
	len int64
	slots EHashSlots
}

func NewEHash(fn HashFn, slotCount int) *EHash {
	return new(EHash).Init(fn, slotCount)
}

func (m *EHash) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *EHash) Init(fn HashFn, slotCount int) *EHash {
	m.fn = fn
	m.slots = make(EHashSlots, slotCount)

	for i, _ := range m.slots {
		m.slots[i].Init()
	}

	return m
}

func (m *EHash) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *EHash) Len() int64 {
	return m.len
}