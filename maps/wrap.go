package maps

import (
	"github.com/fncodr/godbase"
)

type Wrap struct {
	wrapped godbase.Map
}

func (m *Wrap) Clear() {
	m.wrapped.Clear()
}

func (m *Wrap) Cut(start, end godbase.Iter, fn godbase.KVMapFn) godbase.Map {
	return m.wrapped.Cut(start, end, fn)
}

func (m *Wrap) Delete(start, end godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, int) {
	return m.wrapped.Delete(start, end, key, val)
}

func (m *Wrap) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return m.wrapped.Find(start, key, val)
}

func (m *Wrap) First() godbase.Iter {
	return m.wrapped.First()
}

func (m *Wrap) Get(key godbase.Key) (interface{}, bool) {
	return m.wrapped.Get(key)
}

func (m *Wrap) Init(w godbase.Map) *Wrap {
	m.wrapped = w
	return m
}

func (m *Wrap) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	allowMulti bool) (godbase.Iter, bool) {
	return m.wrapped.Insert(start, key, val, allowMulti)
}

func (m *Wrap) Len() int64 {
	return m.wrapped.Len()
}

func (m *Wrap) New() godbase.Map {
	return m.wrapped.New()
}

func (m *Wrap) Set(key godbase.Key, val interface{}) bool {
	return m.wrapped.Set(key, val)
}

func (m *Wrap) String() string {
	return m.wrapped.String()
}
