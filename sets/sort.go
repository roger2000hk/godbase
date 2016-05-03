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

func (self *Sort) Insert(offs int, key godbase.Key) int {
	if i := self.Index(offs, key); i != -1 {
		if self.its[i] == key {
			return i
		}
		
		self.its = append(self.its, nil)
		copy(self.its[i+1:], self.its[i:])
		self.its[i] = key
 		self.len++
		return i
	}

	self.its = append(self.its, key)
	self.len++
	return self.len-1 
}

func (self *Sort) Items() SortIts {
	return self.its
}

func (self *Sort) Len() int64 {
	return int64(self.len)
}
