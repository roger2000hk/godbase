package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"sort"
)

type SortIts []godbase.Key

type Sort struct {
	its SortIts
	len int
}

func (self *Sort) Clone() godbase.Set {
	res := &Sort{its: make(SortIts, self.len), len: self.len}
	copy(res.its, self.its)
	return res
}

func (self *Sort) Delete(offs int, key godbase.Key) int {
	if i := self.Index(offs, key); i != -1 {
		if self.its[i] == key {
			copy(self.its[i:], self.its[i+1:])
			self.len--
			return i
		}
	}

	return -1
}

func (self *Sort) Index(offs int, key godbase.Key) int {
	if i := sort.Search(self.len-offs, func(i int) bool {
		v := self.its[i+offs]
		return key == v || key.Less(v)
	}); i < self.len {
		return i
	}

	return -1
}

func (self *Sort) Insert(offs int, key godbase.Key, multi bool) (int, bool) {
	if i := self.Index(offs, key); i != -1 {
		if self.its[i] == key && !multi {
			return i, false
		}

		isl := len(self.its)
		
 		if isl == self.len  {
			self.its = append(self.its, nil)
			isl++
		}

		copy(self.its[i+1:isl], self.its[i:isl-1])
		self.its[i] = key
 		self.len++
		return i, true
	}

	if self.len < len(self.its) {
		self.its[self.len] = key
	} else {
		self.its = append(self.its, key)
	}

	self.len++
	return self.len-1, true
}

func (self *Sort) Items() SortIts {
	return self.its
}

func (self *Sort) Len() int64 {
	return int64(self.len)

}

func (self *Sort) Resize(len int) *Sort {
	nits := make(SortIts, len)
	copy(nits, self.its)
	self.its = nits
	return self
}
