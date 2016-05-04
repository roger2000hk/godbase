package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
)

type Suffix struct {
	Sort
}

// override to delete all suffixes
func (self *Suffix) Delete(start int, key godbase.Key) int {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Delete(start, godbase.StrKey(sk[i:]))
	}

	return self.Sort.Delete(start, sk)
}

// override to delete all suffixes
func (self *Suffix) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	sk := key.(godbase.StrKey)
	res := int64(0)

	for i := 1; i < len(sk) - 1; i++ {
		_, cnt := self.Sort.DeleteAll(start, end, godbase.StrKey(sk[i:]))
		res += cnt
	}

	i, cnt := self.Sort.DeleteAll(start, end, sk)
	return i, res + cnt 
}

// override to insert all suffixes
func (self *Suffix) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	sk := key.(godbase.StrKey)

	for i := 1; i < len(sk) - 1; i++ {
		self.Sort.Insert(start, godbase.StrKey(sk[i:]), multi)
	}

	return self.Sort.Insert(start, key, multi)
}

