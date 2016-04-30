package idxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"hash"
)


type Reverse struct {
	col cols.Any
	hash hash.Hash64
	recIdHash godbase.UIdHash
	recs maps.Any
}

func NewReverse(c cols.Any, sc int, a *maps.SlabAlloc, ls int) *Reverse {
	i := new(Reverse)
	i.recIdHash.Init()
	
	hashRecId := func(id maps.Key) uint64 {
		return i.recIdHash.Hash(godbase.UId(id.(maps.UIdKey)))
	}

	return i.Init(c, maps.NewHash(maps.NewSlabSlots(sc, hashRecId, a, ls)))
}

func (i *Reverse) Delete(r recs.Any) error {
	id := r.Id()

	if _, ok := i.recs.Delete(nil, nil, maps.UIdKey(id), nil); ok != 1 {
		return recs.NotFound(id)
	}

	return nil
}

func (i *Reverse) Find(start maps.Iter, key maps.Key, val interface{}) (maps.Iter, bool) {
	return i.recs.Find(start, key, val)	
}

func (i *Reverse) Init(c cols.Any, rs maps.Any) *Reverse {
	i.col = c
	i.recs = rs
	return i
}

func (i *Reverse) Insert(r recs.Any) (recs.Any, error) {
	i.recs.Set(maps.UIdKey(r.Id()), r.Get(i.col))
	return r, nil
}
