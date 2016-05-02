package sets

// package godbase/sets implements hashed and sorted sets based on slices

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"sort"
)

type Sort []godbase.Key


func (self Sort) Index(key godbase.Key, offs int) int {
	return sort.Search(len(self)-offs, func(i int) bool {
		v := self[i+offs]
		return key == v || key.Less(v)
	})
}

func (self Sort) Delete(key godbase.Key) (godbase.Set, bool) {
	l := len(self)
	if i := self.Index(key, 0); i < l {
		if self[i] == key {
			ns := make(Sort, l-1)
			copy(ns, self[:i])
			copy(ns[i:], self[i+1:])
			return ns, true
		}
	}

	return self, false
}

func (self Sort) HasKey(key godbase.Key) bool {
	if i := self.Index(key, 0); i < len(self) {
		return self[i] == key
	}

	return false
}

func (self Sort) Insert(key godbase.Key) (godbase.Set, bool) {
	l := len(self)

	if i := self.Index(key, 0); i < l {
		if self[i] == key {
			return self[i:], false
		}
		
		ns := make(Sort, l+1)
		copy(ns, self[:i])
		ns[i] = key
		copy(ns[i+1:], self[i:])
		return ns, true
	}

	return append(self, key), true 
}

func (self Sort) Len() int {
	return len(self)
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
