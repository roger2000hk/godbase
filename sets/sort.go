package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/utils"
	"sort"
)

type SortElems []godbase.Key

type Sort struct {
	elems SortElems
	len int
}

func NewSort() *Sort {
	return new(Sort)
}

func (self *Sort) Capac() int {
	return len(self.elems)
}

func (self *Sort) Clone() godbase.Set {
	res := &Sort{elems: make(SortElems, self.len), len: self.len}
	copy(res.elems, self.elems)
	return res
}

func (self *Sort) Delete(start int, key godbase.Key) int {
	if i := self.index(start, key); i != self.len {
		if self.elems[i] == key {
			copy(self.elems[i:], self.elems[i+1:])
			self.len--
			return i
		}
	}

	return -1
}

func (self *Sort) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	i := start

	if end == 0 {
		end = self.len
	}

	if key != nil {
		i = self.index(start, key)
	}

	var j int
	cnt := 0
	
	for j = i; j < end && self.elems[j] == key; j++ {
		cnt++
	}

	if j > i {
		copy(self.elems[i:], self.elems[j:])
		self.len -= cnt
	}
	
	return i, int64(cnt)
}

func (self *Sort) First(start int, key godbase.Key) (int, bool) {
	i := self.index(start, key)

	if i < self.len {
		if self.elems[i] == key {
			return i, true
		}
		
		if i > 0 {
			return i-1, false
		}
	}

	return i, false
}

func (self *Sort) Get(_ godbase.Key, i int) godbase.Key {
	return self.elems[i]
}

func (self *Sort) Last(start, end int, key godbase.Key) (int, bool) {
	i := start

	if end == 0 {
		end = self.len
	}

	if key != nil {
		i = self.index(start, key)
	}

	var j int
	res := false

	for j = i; j < end && self.elems[j] == key; j++ {
		res = true
	}
		
	return j-1, res
}

func (self *Sort) Load(start int, keys...godbase.Key) {
	i := self.index(start, keys[0])
	ksl := len(keys)
	nvs := make(SortElems, self.len + ksl)
	copy(nvs, self.elems[:i])
	copy(nvs[:i], keys)
	copy(nvs[:i+ksl], self.elems[i:])
	self.len += ksl
}

func (self *Sort) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	if i := self.index(start, key); i < self.len {
		if self.elems[i] == key && !multi {
			return i, false
		}

		isl := len(self.elems)
		
 		if isl == self.len  {
			self.elems = append(self.elems, nil)
			isl++
		}

		copy(self.elems[i+1:isl], self.elems[i:isl-1])
		self.elems[i] = key
 		self.len++
		return i, true
	}

	if self.len < len(self.elems) {
		self.elems[self.len] = key
	} else {
		self.elems = append(self.elems, key)
	}

	self.len++
	return self.len-1, true
}

func (self *Sort) Len() int64 {
	return int64(self.len)

}

func (self *Sort) Resize(s int) *Sort {
	nelems := make(SortElems, s)
	copy(nelems, self.elems[:utils.Min(s, len(self.elems))])
	self.elems = nelems
	return self
}

func (self *Sort) While(fn godbase.IKTestFn) bool {
	for i, k := range self.elems {
		if !fn(i, k) {
			return false
		}
	}
	
	return true
}

func (self *Sort) index(start int, key godbase.Key) int {
	return sort.Search(self.len-start, func(i int) bool {
		v := self.elems[i+start]
		return key == v || key.Less(v)
	})
}
