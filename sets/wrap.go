package sets

import (
	"github.com/fncodr/godbase"
)

type Wrap struct {
	wrapped godbase.Set
}

func (self *Wrap) Clone() godbase.Set {
	return self.wrapped.Clone()
}

func (self *Wrap) Delete(start int, key godbase.Key) int {
	return self.wrapped.Delete(start, key)
}

func (self *Wrap) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	return self.wrapped.DeleteAll(start, end, key)
}

func (self *Wrap) First(start int, key godbase.Key) int {
	return self.wrapped.First(start, key)
}

func (self *Wrap) Get(key godbase.Key, i int) godbase.Key {
	return self.wrapped.Get(key, i)
}

func (self *Wrap) Last(start, end int, key godbase.Key) int {
	return self.wrapped.Last(start, end, key)
}

func (self *Wrap) Init(s godbase.Set) *Wrap {
	self.wrapped = s
	return self
}

func (self *Wrap) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	return self.wrapped.Insert(start, key, multi)
}

func (self *Wrap) Len() int64 {
	return self.wrapped.Len()
}

func (self *Wrap) While(fn godbase.IKTestFn) bool {
	return self.wrapped.While(fn)
}
