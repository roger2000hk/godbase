package maps

type Suffix struct {
	Wrap
}

func NewSuffix(m Any) *Suffix {
	return &Suffix{Wrap: Wrap{wrapped: m}}
}

func (m *Suffix) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	sk := key.(StringKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := m.wrapped.Delete(start, end, StringKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := m.wrapped.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

func (m *Suffix) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	sk := key.(StringKey)

	for i := 1; i < len(sk) - 1; i++ {
		m.wrapped.Insert(start, StringKey(sk[i:]), val, allowMulti)
	}

	return m.wrapped.Insert(start, key, val, allowMulti)
}

