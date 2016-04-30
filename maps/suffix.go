package maps

import (
	//"fmt"
	"github.com/fncodr/godbase"
)

type Suffix struct {
	Wrap
}

func NewSuffix(m godbase.Map) *Suffix {
	res := new(Suffix)
	res.Init(m)
	return res
}

// override to delete all suffixes
func (m *Suffix) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	sk := key.(godbase.StringKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := m.wrapped.Delete(start, end, godbase.StringKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := m.wrapped.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

// override to insert all suffixes
func (m *Suffix) Insert(start godbase.Iter, key godbase.Key, val interface{}, allowMulti bool) (godbase.Iter, bool) {
	sk := key.(godbase.StringKey)

	for i := 1; i < len(sk) - 1; i++ {
		m.wrapped.Insert(start, godbase.StringKey(sk[i:]), val, allowMulti)
	}

	return m.wrapped.Insert(start, key, val, allowMulti)
}

