package sets

import (
	//"fmt"
	"github.com/fncodr/godbase"
	"sort"
)

type Sort []godbase.Key

func (self Sort) Delete(offs int, key godbase.Key) (godbase.Set, int) {
	if i := self.Index(offs, key); i != -1 {
		if self[i] == key {
			copy(self[i:], self[i+1:])
			return self[:len(self)-1], i
		}
	}

	return self, 1
}

func (self Sort) Index(offs int, key godbase.Key) int {
	if i := sort.Search(len(self)-offs, func(i int) bool {
		v := self[i+offs]
		return key == v || key.Less(v)
	}); i < len(self) {
		return i
	}

	return -1
}

func (self Sort) Insert(offs int, key godbase.Key) (godbase.Set, int) {
	if i := self.Index(offs, key); i != -1 {
		if self[i] == key {
			return self, -1
		}
		
		self = append(self, nil)
		copy(self[i+1:], self[i:])
		self[i] = key
		return self, i
	}

	return append(self, key), len(self) 
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
