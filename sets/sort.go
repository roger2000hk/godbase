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

func (self *Sort) Delete(start int, key godbase.Key) int {
	if i := self.index(start, key); i != self.len {
		if self.its[i] == key {
			copy(self.its[i:], self.its[i+1:])
			self.len--
			return i
		}
	}

	return -1
}

func (self *Sort) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	i := start

	if end == -1 {
		end = self.len
	}

	if key != nil {
		i = self.index(start, key)
	}

	var j int
	cnt := 0
	
	for j = i; j < end && self.its[j] == key; j++ {
		cnt++
	}

	if j > i {
		copy(self.its[i:], self.its[j:])
		self.len -= cnt
	}
	
	return i, int64(cnt)
}

func (self *Sort) First(start int, key godbase.Key) int {
	if i := self.index(start, key); i < self.len && self.its[i] == key {
		return i
	}

	return -1
}

func (self *Sort) Last(start, end int, key godbase.Key) int {
	i := start

	if end == -1 {
		end = self.len
	}

	if key != nil {
		i = self.index(start, key)
	}

	var j int

	for j = i; j < end && self.its[j] == key; j++ {
	}
		
	return j-1
}

func (self *Sort) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	if i := self.index(start, key); i < self.len {
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

func (self *Sort) index(start int, key godbase.Key) int {
	return sort.Search(self.len-start, func(i int) bool {
		v := self.its[i+start]
		return key == v || key.Less(v)
	})
}
