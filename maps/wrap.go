package maps

import (
	"github.com/fncodr/godbase"
)

type Wrap struct {
	wrapped godbase.Map
}

func (self *Wrap) Clear() {
	self.wrapped.Clear()
}

func (self *Wrap) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	return self.wrapped.Cut(start, end, fn)
}

func (self *Wrap) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	return self.wrapped.Delete(start, end, key, val)
}

func (self *Wrap) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return self.wrapped.Find(start, key, val)
}

func (self *Wrap) First() godbase.Iter {
	return self.wrapped.First()
}

func (self *Wrap) Get(key godbase.Key) (interface{}, bool) {
	return self.wrapped.Get(key)
}

func (self *Wrap) Init(w godbase.Map) *Wrap {
	self.wrapped = w
	return self
}

func (self *Wrap) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	multi bool) (godbase.Iter, bool) {
	return self.wrapped.Insert(start, key, val, multi)
}

func (self *Wrap) Len() int64 {
	return self.wrapped.Len()
}

func (self *Wrap) New() godbase.Map {
	return self.wrapped.New()
}

func (self *Wrap) Set(key godbase.Key, val interface{}) bool {
	return self.wrapped.Set(key, val)
}

func (self *Wrap) String() string {
	return self.wrapped.String()
}
