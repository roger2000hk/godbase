package maps

type Any interface {
	Delete(key Cmp, val interface{}) int
	Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool)
	Len() int64
}

type Cmp interface {
	Less(Cmp) bool
}

type Map map[Cmp]interface{}

func NewMap() Map {
	return make(Map)
}

func (m Map) Delete(key Cmp, val interface{}) int {
	delete(m, key)
	return 1
}

func (m Map) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	if prev, ok := m[key]; ok {
		return prev, false
	}
	
	m[key] = val
	return val, true
}

func (m Map) Len() int64 {
	return int64(len(m))
}

