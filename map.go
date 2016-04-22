package godbase

type Map interface {
	Delete(key Cmp, val interface{}) int
	Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool)
	Len() int64
}

type Cmp interface {
	Less(Cmp) bool
}

type MapMap map[Cmp]interface{}

func NewMapMap() MapMap {
	return make(MapMap)
}

func (m MapMap) Delete(key Cmp, val interface{}) int {
	delete(m, key)
	return 1
}

func (m MapMap) Insert(key Cmp, val interface{}, allowMulti bool) (interface{}, bool) {
	if prev, ok := m[key]; ok {
		return prev, false
	}
	
	m[key] = val
	return val, true
}

func (m MapMap) Len() int64 {
	return int64(len(m))
}

