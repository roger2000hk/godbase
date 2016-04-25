package maps

type Wrap struct {
	wrapped Any
}

func (m *Wrap) Cut(start, end Iter, fn TestFn) Any {
	return m.wrapped.Cut(start, end, fn)
}

func (m *Wrap) Delete(start, end Iter, key Key, val interface{}) (Iter, int) {
	return m.wrapped.Delete(start, end, key, val)
}

func (m *Wrap) Find(start Iter, key Key, val interface{}) (Iter, bool) {
	return m.wrapped.Find(start, key, val)
}

func (m *Wrap) Insert(start Iter, key Key, val interface{}, allowMulti bool) (Iter, bool) {
	return m.wrapped.Insert(start, key, val, allowMulti)
}

func (m *Wrap) Len() int64 {
	return m.wrapped.Len()
}

func (m *Wrap) String() string {
	return m.wrapped.String()
}