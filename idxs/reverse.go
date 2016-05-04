package idxs

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"hash"
)


type Reverse struct {
	Basic
	col godbase.Col
	hash hash.Hash64
	recIdHash godbase.UIdHash
	recs godbase.Map
}

func NewReverse(n string, c godbase.Col, sc int, a *maps.SlabAlloc, ls int) *Reverse {
	i := new(Reverse)
	i.recIdHash.Init()
	
	hashRecId := func(id godbase.Key) uint64 {
		return i.recIdHash.Hash(godbase.UId(id.(godbase.UIdKey)))
	}

	return i.Init(n, c, maps.NewSlabHash(sc, hashRecId, a, ls))
}

func (self *Reverse) AddToTbl(tbl godbase.Tbl)  {
	tbl.AddIdx(self)
}

func (i *Reverse) Delete(start godbase.Iter, r godbase.Rec) error {
	return i.Drop(start, r)
}

func (i *Reverse) Drop(start godbase.Iter, r godbase.Rec) error {
	id := r.Id()

	if _, ok := i.recs.Delete(start, nil, godbase.UIdKey(id), nil); ok != 1 {
		return recs.NotFound(id)
	}

	return nil
}

func (i *Reverse) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return i.recs.Find(start, key, val)	
}

func (i *Reverse) Init(n string, c godbase.Col, rs godbase.Map) *Reverse {
	i.Basic.Init(n)
	i.col = c
	i.recs = rs
	return i
}

func (i *Reverse) Insert(start godbase.Iter, r godbase.Rec) (godbase.Iter, error) {
	k, v := godbase.UIdKey(r.Id()), r.Get(i.col)
	res, ok := i.recs.Insert(start, k, v, false)

	if !ok && !i.col.Eq(res.Val(), v) {
		return nil, &DupKey{key: k}
	}

	return res, nil
}

func (i *Reverse) Key(vs...interface{}) godbase.Key {
	if len(vs) > 1 {
		panic(fmt.Sprintf("invalid suffix key: %v", vs))
	}

	return godbase.UIdKey(vs[0].(godbase.Rec).Id())
}

func (i *Reverse) Load(rec godbase.Rec) (godbase.Rec, error) {
	i.recs.Set(godbase.UIdKey(rec.Id()), rec.Get(i.col))
	return rec, nil
}

func (i *Reverse) RecKey(rec godbase.Rec) godbase.Key {
	return godbase.UIdKey(rec.Id())
}

