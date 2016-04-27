package tbls

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"io"
)

type Any interface {
	defs.Any
	Add(cols.Any) cols.Any
	Col(n string) cols.Any
	Cols() ColIter
	Read(rec recs.Any, r io.Reader) (recs.Any, error)
	Upsert(rec recs.Any) recs.Any
	Write(recs.Any, io.Writer) error
}

type Basic struct {
	defs.Basic
	cols maps.Skip
	recIdHash recs.IdHash
	recs maps.Hash
}

type RecAlloc *maps.SkipAlloc
type ColIter maps.Iter
type RecIterIter maps.Iter

func New(n string, rsc int, ra RecAlloc, rls int) Any {
	return new(Basic).Init(n, rsc, ra, rls)
}

func (t *Basic) Col(n string) cols.Any {
	if c, ok := t.cols.Get(maps.StringKey(n)); ok {
		return c.(cols.Any)
	}
	
	panic(fmt.Sprintf("col not found: %v", n))
}

func (t *Basic) Cols() ColIter {
	return t.cols.First()
}

func (t *Basic) Add(c cols.Any) cols.Any {
	return t.cols.Set(maps.StringKey(c.Name()), c).(cols.Any)
}

func (t *Basic) Init(n string, rsc int, ra RecAlloc, rls int) *Basic {
	t.Basic.Init(n)
	t.cols.Init(nil, 1)
	t.recIdHash.Init()

	hashRecId := func(_id maps.Key) uint64 {
		id := _id.(recs.Id)
		return t.recIdHash.Hash(id)
	}

	t.recs.Init(maps.NewSkipSlots(rsc, hashRecId, ra, rls))
	t.Add(recs.CreatedAtCol())
	t.Add(recs.IdCol())
	return t
}

func (t *Basic) Read(rec recs.Any, r io.Reader) (recs.Any, error) {
	var s recs.Size

	if err := godbase.Read(&s, r); err != nil {
		return nil, err
	}
	
	for i := recs.Size(0); i < s; i++ {
		var n string
		var err error

		if n, err = cols.ReadName(r); err != nil {
			return nil, err
		}

		c := t.Col(n)
		var v interface{}
		if v, err = cols.Read(c, r); err != nil {
			return nil, err
		}

		rec.Set(c, v)
	}

	return rec, nil
}

func (t *Basic) Upsert(rec recs.Any) recs.Any {
	id := rec.Id()
	rr := rec.New()
	
	for i := r.cols.First(); i.Valid(); i = i.Next() {
		c := i.Val().(cols.Any)
		if v, ok := rec.Get(c); ok {
			rr.Set(c, v)
		}
	}
	
	t.recs.Set(id, rr)
	return rec
}

func (t *Basic) Write(rec recs.Any, w io.Writer) error {
	s := recs.Size(rec.Len())

	if err := godbase.Write(&s, w); err != nil {
		return err
	}

	for i := rec.Iter(); i.Valid(); i=i.Next() {
		if err := cols.Write(i.Key().(cols.Any), i.Val(), w); err != nil {
			return err
		}
	}

	return nil
}
