package godbase

type EHashSlots []*ESkipMap

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

	if s := m.slots[i]; s != nil {
		res := s.Delete(key, val)
		m.len -= int64(res)
		return res
	}

	return 0
}

func (m *EHashMap) Init(fn HashFn, slotCount int) *EHashMap {
	m.fn = fn
	m.slots = make(EHashSlots, slotCount)
	return m
}

func (m *EHashMap) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	res, ok := m.getSlot(key).Insert(key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *EHashMap) Len() int64 {
	return m.len
}

func (m *EHashMap) getSlot(key Cmp) *ESkipMap {
	i := m.fn(key) % uint64(len(m.slots))
	s := m.slots[i]

	if s == nil {
		s = NewESkipMap()
		m.slots[i] = s
	}

	return s
}
