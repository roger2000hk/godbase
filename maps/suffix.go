package maps

import (
	//"fmt"
	"github.com/fncodr/godbase"
)

type Suffix struct {
	Sort
}

func NewSuffix(a *SlabAlloc, ls int) *Suffix {
	res := new(Suffix)
	res.Sort.Init(a, ls)
	return res
}

// override to delete all suffixes
func (m *Suffix) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	sk := key.(godbase.StringKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := m.Sort.Delete(start, end, godbase.StringKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := m.Sort.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

// override to insert all suffixes
func (m *Suffix) Insert(start godbase.Iter, key godbase.Key, val interface{}, multi bool) (godbase.Iter, bool) {
	sk := key.(godbase.StringKey)

	for i := 1; i < len(sk) - 1; i++ {
		m.Sort.Insert(start, godbase.StringKey(sk[i:]), val, multi)
	}

	return m.Sort.Insert(start, key, val, multi)
}

