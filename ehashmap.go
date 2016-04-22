package godbase

type EHashSlots []ESkipMap

type EHashMap struct {
	fn HashFn
	len int64
	slots EHashSlots
}

func NewEHashMap(fn HashFn, slotCount int) *EHashMap {
	return new(EHashMap).Init(fn, slotCount)
}

func (m *EHashMap) Delete(key Cmp, val interface{}) int {
	i := m.fn(key) % uint64(len(m.slots))
	res := m.slots[i].Delete(key, val)
	m.len -= int64(res)
	return res
}

func (m *EHashMap) Init(fn HashFn, slotCount int) *EHashMap {
	m.fn = fn
	m.slots = make(EHashSlots, slotCount)

	for i, _ := range m.slots {
		m.slots[i].Init()
	}

	return m
}

func (m *EHashMap) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	i := m.fn(key) % uint64(len(m.slots))
	res, ok := m.slots[i].Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *EHashMap) Len() int64 {
	return m.len
}
