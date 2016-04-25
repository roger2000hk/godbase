package maps

type Slots interface {	
	Get(key Cmp, create bool) Any
}

type Hash struct {
	len int64
	slots Slots
}

func NewHash(slots Slots) *Hash {
	return new(Hash).Init(slots)
}

func (m *Hash) Cut(start, end Iter, fn TestFn) Any {
	return m.slots.Get(start.Key(), true).Cut(start, end, fn)
}

func (m *Hash) Delete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	res, cnt := m.slots.Get(key, true).Delete(start, end, key, val)
	m.len -= int64(cnt)
	return res, cnt
}

func (m *Hash) Find(start Iter, key Cmp, val interface{}) (Iter, bool) {
	res, ok := m.slots.Get(key, true).Find(start, key, val)
	return res, ok
}

func (m *Hash) Init(slots Slots) *Hash {
	m.slots = slots
	return m
}

func (m *Hash) Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {	
	res, ok := m.slots.Get(key, true).Insert(start, key, val, allowMulti)

	if ok {
		m.len++
	}

	return res, ok
}

func (m *Hash) Len() int64 {
	return m.len
}

type SkipSlots struct {
	alloc *SkipAlloc
	fn HashFn
	levels int
	slots []Skip
}

func NewSkipSlots(count int, fn HashFn, alloc *SkipAlloc, levels int) *SkipSlots {
	ss := new(SkipSlots)
	ss.alloc = alloc
	ss.fn = fn
	ss.levels = levels
	ss.slots = make([]Skip, count)
	return ss
}

func (ss SkipSlots) Get(key Cmp, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := &ss.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init(ss.alloc, ss.levels)
	}

	return nil
}

type ESkipSlots struct {
	fn HashFn
	slots []ESkip
}

func NewESkipSlots(count int, fn HashFn) *ESkipSlots {
	ss := new(ESkipSlots)
	ss.fn = fn
	ss.slots = make([]ESkip, count)
	return ss
}

func (ss ESkipSlots) Get(key Cmp, create bool) Any {
	i := ss.fn(key) % uint64(len(ss.slots))
	s := &ss.slots[i]

	if s.isInit {
		return s
	}

	if create {
		return s.Init()
	}

	return nil
}
