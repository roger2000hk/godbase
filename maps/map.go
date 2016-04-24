package maps

// Map is mostly meant as a reference for performance comparisons,
// it only supports enough of the api to run the basic tests

type Map map[Cmp]interface{}

func NewMap() Map {
	return make(Map)
}

func (m Map) Cut(start, end Iter, fn TestFn) Any {
	panic("Map doesn't support iters")
}

func (m Map) Delete(start, end Iter, key Cmp, val interface{}) (Iter, int) {
	if start != nil || end != nil {
		panic("Map doesn't support iters")
	}

	if val != nil {
		panic("Map doesn't support multi")
	}

	delete(m, key)
	return nil, 1
}

func (m Map) Find(start Iter, key Cmp, val interface{}) (Iter, bool) {
	panic("Map doesn't support iters")
}

func (m Map) Insert(start Iter, key Cmp, val interface{}, allowMulti bool) (Iter, bool) {
	if start != nil {
		panic("Map doesn't support iters")
	}

	if allowMulti {
		panic("Map doesn't support multi")
	}

	if _, ok := m[key]; ok {
		return nil, false
	}
	
	m[key] = val
	return nil, true
}

func (m Map) Len() int64 {
	return int64(len(m))
}
