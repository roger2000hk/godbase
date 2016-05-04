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
func (self *Suffix) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	sk := key.(godbase.StrKey)
	cnt := 0

	for i := 1; i < len(sk) - 1; i++ {
		_, sc := self.Sort.Delete(start, end, godbase.StrKey(sk[i:]), val)
		cnt += sc
	}

	res, sc := self.Sort.Delete(start, end, sk, val)
	cnt += sc
	return res, cnt
}

// override to insert all suffixes
func (self *Suffix) Insert(start godbase.Iter, key godbase.Key, val interface{}, multi bool) (godbase.Iter, bool) {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Insert(start, godbase.StrKey(sk[i:]), val, multi)
	}

	return self.Sort.Insert(start, key, val, multi)
}

