package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"sort"
)

type Sort []godbase.Key

func (self Sort) Delete(offs int, key godbase.Key) (godbase.Set, bool) {
	l := len(self)
	if i := self.Index(offs, key); i < l {
		if self[i] == key {
			copy(self[i:], self[i+1:])
			return self, true
		}
	}

	return self, false
}

func (self Sort) HasKey(offs int, key godbase.Key) bool {
	if i := self.Index(offs, key); i < len(self) {
		return self[i] == key
	}

	return false
}

func (self Sort) Index(offs int, key godbase.Key) int {
	return sort.Search(len(self)-offs, func(i int) bool {
		v := self[i+offs]
		return key == v || key.Less(v)
	})
}

func (self Sort) Insert(offs int, key godbase.Key) (godbase.Set, bool) {
	l := len(self)

	if i := self.Index(offs, key); i < l {
		if self[i] == key {
			return self[i:], false
		}
		
		self = append(self, nil)
		copy(self[i+1:], self[i:])
		self[i] = key
		return self, true
	}

	return append(self, key), true 
}

func (self Sort) Len() int64 {
	return int64(len(self))
}

func (self Sort) Key() godbase.Key {
	if len(self) == 0 {
		return nil
	}

	return self[0]
}

func (self Sort) Next() godbase.Iter {
	if len(self) == 0 {
		return self
	}
	
	return self[1:]
}

func (self Sort) Val() interface{} {
	if len(self) == 0 {
		return nil
	}

	return self[0]
}

func (self Sort) Valid() bool {
	return len(self) != 0
}
