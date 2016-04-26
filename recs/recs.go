package recs

import (
	"fmt"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
)

type Rec interface {
	Int64(*cols.Int64) int64
	SetInt64(*cols.Int64, int64) int64
	SetUInt64(*cols.UInt64, uint64) uint64
	UInt64(*cols.UInt64) uint64
}

type BasicRec maps.Skip

func NewRec(alloc *maps.SkipAlloc, levels int) Rec {
	return (*BasicRec)(maps.NewSkip(alloc, levels))
}

func (r *BasicRec) Int64(c *cols.Int64) int64 {
	if v, ok := r.AsMap().Get(nil, c); ok {
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
	if v, ok := r.AsMap().Get(nil, c); ok {
		return v.(uint64)
	}

	panic(fmt.Sprintf("field not found: %v", c.Name()))
}

func (r *BasicRec) AsMap() *maps.Skip {
	return (*maps.Skip)(r)
}
