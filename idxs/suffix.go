package idxs

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
)

type Suffix struct {
	Basic
	cols []*cols.StringCol
	recs maps.Suffix
	unique bool
}

func NewSuffix(n string, cs []*cols.StringCol, u bool, a *maps.SlabAlloc, ls int) *Suffix {
	return new(Suffix).Init(n, cs, u, a, ls)
}

func (self *Suffix) AddToTbl(tbl godbase.Tbl)  {
	tbl.AddIdx(self)
}

func (i *Suffix) Delete(start godbase.Iter, r godbase.Rec) error {
	return i.Drop(start, r)
}

func (i *Suffix) Drop(start godbase.Iter, r godbase.Rec) error {
	id := r.Id()

	for _, col := range i.cols {
		if v, ok := r.Find(col); ok {
			if _, ok := i.recs.Delete(start, nil, godbase.StringKey(v.(string)), id); 
			ok != 1 {
				return recs.NotFound(id)
			}
		}
	}

	return nil
}

func (self *Suffix) Find(start godbase.Iter, _key godbase.Key, val interface{}) (godbase.Iter, bool) {
	res := start
	key := _key.(godbase.StringsKey)
	
	for i, _ := range self.cols {
		if len(key) < i+1 {
			break
		}

		var ok bool
		if res, ok = self.recs.Find(start, godbase.StringKey(key[i]), val); ok {
			return res, true
		}
	}

	return res, false
}

func (i *Suffix) Init(n string, cs []*cols.StringCol, u bool, a *maps.SlabAlloc, ls int) *Suffix {
	i.Basic.Init(n)
	i.recs.Init(a, ls)
	i.cols = cs
	i.unique = u
	return i
}

func (i *Suffix) Insert(start godbase.Iter, r godbase.Rec) (godbase.Iter, error) {
	id := r.Id()
	res := start

	for _, col := range i.cols {
		if v, ok := r.Find(col); ok {
			k := godbase.StringKey(v.(string))
			var ok bool
			res, ok = i.recs.Insert(start, k, id, !i.unique)
			
			if !ok && !col.Eq(res.Val(), v) {
				return nil, &DupKey{key: k}
			}
		}
	}
	
	return res, nil
}

func (i *Suffix) Key(vs...interface{}) godbase.Key {
	res := make([]string, len(vs))

	for i, v := range vs {
		res[i] = v.(string)
	}

	return godbase.StringsKey(res)
}

func (i *Suffix) Load(rec godbase.Rec) (godbase.Rec, error) {
	for _, col := range i.cols {
		if v, ok := rec.Find(col); ok {
			i.recs.Set(godbase.StringKey(v.(string)), rec.Id())
		}
	}

	return rec, nil
}

func (i *Suffix) RecKey(rec godbase.Rec) godbase.Key {
	var res []string

	for _, c := range i.cols {
		if v, ok := rec.Find(c); ok {
			res = append(res, v.(string))
		}
	}

	return godbase.StringsKey(res)
}

