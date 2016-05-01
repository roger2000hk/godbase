package idxs

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
)

type Suffix struct {
	defs.Basic
	col *cols.StringCol
	recs maps.Suffix
	unique bool
}

func NewSuffix(n string, c *cols.StringCol, u bool, a *maps.SlabAlloc, ls int) *Suffix {
	return new(Suffix).Init(n, c, u, a, ls)
}

func (i *Suffix) Delete(start godbase.Iter, r godbase.Rec) error {
	return i.Drop(start, r)
}

func (i *Suffix) Drop(start godbase.Iter, r godbase.Rec) error {
	id := r.Id()

	if v, ok := r.Find(i.col); ok {
		if _, ok := i.recs.Delete(start, nil, godbase.StringKey(v.(string)), id); ok != 1 {
			return recs.NotFound(id)
		}
	}

	return nil
}

func (i *Suffix) Find(start godbase.Iter, key godbase.Key, val interface{}) (godbase.Iter, bool) {
	return i.recs.Find(start, key, val)	
}

func (i *Suffix) Init(n string, c *cols.StringCol, u bool, a *maps.SlabAlloc, ls int) *Suffix {
	i.Basic.Init(n)
	i.recs.Init(a, ls)
	i.col = c
	i.unique = u
	return i
}

func (i *Suffix) Insert(start godbase.Iter, r godbase.Rec) (godbase.Iter, error) {
	id := r.Id()

	if v, ok := r.Find(i.col); ok {
		k := godbase.StringKey(v.(string))
		res, ok := i.recs.Insert(start, k, id, i.unique)

		if !ok && !i.col.Eq(res.Val(), v) {
			return nil, &DupKey{key: k}
		}

		return res, nil
	}
	
	return start, nil
}

func (i *Suffix) Key(vs...interface{}) godbase.Key {
	if len(vs) > 1 {
		panic(fmt.Sprintf("invalid suffix key: %v", vs))
	}

	return godbase.StringKey(vs[0].(string))
}

func (i *Suffix) Load(rec godbase.Rec) (godbase.Rec, error) {
	if v, ok := rec.Find(i.col); ok {
		i.recs.Set(godbase.StringKey(v.(string)), rec.Id())
	}
	
	return rec, nil
}

func (i *Suffix) RecKey(rec godbase.Rec) godbase.Key {
	return godbase.StringKey(recs.String(rec, i.col))
}

