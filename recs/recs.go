package recs

import (
	"fmt"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
)

type Rec interface {
	Delete(cols.Any) bool
	Get(cols.Any) (interface{}, bool)
	Int64(*cols.Int64) int64
	SetInt64(*cols.Int64, int64) int64
	SetUInt64(*cols.UInt64, uint64) uint64
	UInt64(*cols.UInt64) uint64
}

type BasicRec maps.Skip

func NewRec(alloc *maps.SkipAlloc, levels int) Rec {
	return new(BasicRec).Init(alloc, levels)
}

func (r *BasicRec) AsMap() *maps.Skip {
	return (*maps.Skip)(r)
}

func (r *BasicRec) Delete(c cols.Any) bool {
	_, cnt := r.AsMap().Delete(nil, nil, c, nil)
	return cnt == 1
}

func (r *BasicRec) Get(c cols.Any) (interface{}, bool) {
	if v, ok := r.AsMap().Get(nil, c); ok {
		return v, true
	}
	
	return nil, false
}

func (r *BasicRec) Init(alloc *maps.SkipAlloc, levels int) *BasicRec {
	r.AsMap().Init(alloc, levels)
	return r
}

func (r *BasicRec) Int64(c *cols.Int64) int64 {
	if v, ok := r.Get(c); ok {
		return v.(int64)
	}

	panic(fmt.Sprintf("field not found: %v", c.Name()))
}

func (r *BasicRec) SetInt64(c *cols.Int64, v int64) int64 {
	r.AsMap().Insert(nil, c, v, false)
	return v
}

func (r *BasicRec) SetUInt64(c *cols.UInt64, v uint64) uint64 {
	r.AsMap().Insert(nil, c, v, false)
	return v
}

func (r *BasicRec) UInt64(c *cols.UInt64) uint64 {
	if v, ok := r.Get(c); ok {
		return v.(uint64)
	}

	panic(fmt.Sprintf("field not found: %v", c.Name()))
}
