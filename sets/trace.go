package sets

import (
	"github.com/fncodr/godbase"
	"log"
)

// embedding Wrap gives you default delegation to wrapped set

type Trace struct {
	Wrap

	// we're adding an id for logging
	id string
}

func NewTrace(s godbase.Set, id string) *Trace {
	res := new(Trace)
	res.Init(s)
	res.id = id
	return res
}

// override to log actions before updating wrapped map

func (self *Trace) Delete(start int, key godbase.Key) int {
	log.Printf("%v.Delete %v: '%v'", self.id, start, key)
	return self.wrapped.Delete(start, key)
}

func (self *Trace) DeleteAll(start, end int, key godbase.Key) (int, int64) {
	log.Printf("%v.DeleteAll %v/%v: '%v'", self.id, start, end, key)
	return self.wrapped.DeleteAll(start, end, key)
}

func (self *Trace) Insert(start int, key godbase.Key, multi bool) (int, bool) {
	log.Printf("%v.Insert/%v %v: '%v'", self.id, multi, start, key)
	return self.wrapped.Insert(start, key, multi)
}
