package idxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"hash"
)


type Reverse struct {
	col godbase.Col
	hash hash.Hash64
	recIdHash godbase.UIdHash
	recs godbase.Map
}

func NewReverse(c godbase.Col, sc int, a *maps.SlabAlloc, ls int) *Reverse {
	i := new(Reverse)
	i.recIdHash.Init()
	
	hashRecId := func(id godbase.Key) uint64 {
		return i.recIdHash.Hash(godbase.UId(id.(godbase.UIdKey)))
	}

	return i.Init(c, maps.NewHash(maps.NewSlabSlots(sc, hashRecId, a, ls)))
}

func (i *Reverse) Delete(r godbase.Rec) error {
	id := r.Id()

	if _, ok := i.recs.Delete(nil, nil, godbase.UIdKey(id), nil); ok != 1 {
		return recs.NotFound(id)
	}

	return nil
}

func (i *Reverse) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return i.recs.Find(start, key, val)	
}

func (i *Reverse) Init(c godbase.Col, rs godbase.Map) *Reverse {
	i.col = c
	i.recs = rs
	return i
}

func (i *Reverse) Insert(r godbase.Rec) (godbase.Rec, error) {
	i.recs.Set(godbase.UIdKey(r.Id()), r.Get(i.col))
	return r, nil
}
