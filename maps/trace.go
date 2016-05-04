package maps

import (
	"github.com/fncodr/godbase"
	"log"
)

// embedding Wrap gives you default delegation to wrapped map

type Trace struct {
	Wrap

	// we're adding an id for logging
	id string
}

func NewTrace(m godbase.Map, id string) *Trace {
	res := new(Trace)
	res.Init(m)
	res.id = id
	return res
}

// override to log actions before updating wrapped map

func (self *Trace) Delete(start, end godbase.Iter, key godbase.Key, 
	val interface{}) (godbase.Iter, int) {
	log.Printf("%v.delete '%v': '%v'", self.id, key, val)
	return self.wrapped.Delete(start, end, key, val)
}

func (self *Trace) Insert(start godbase.Iter, key godbase.Key, val interface{}, 
	multi bool) (godbase.Iter, bool) {
	log.Printf("%v.insert/%v '%v': '%v'", self.id, multi, key, val)
	return self.wrapped.Insert(start, key, val, multi)
}
