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

type ColIter maps.ValIter

type Any interface {
	defs.Any
	Add(cols.Any) cols.Any
	Col(n string) cols.Any
	Cols() ColIter
	Write(recs.Any, io.Writer) error
}

type Basic struct {
	cols maps.Skip
	defs.Basic
}

func New(n string) Any {
	return new(Basic).Init(n)
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

func (t *Basic) Init(n string) *Basic {
	t.Basic.Init(n)
	t.cols.Init(nil, 1)
	t.Add(recs.IdCol())
	return t
}

func (t *Basic) Write(r recs.Any, w io.Writer) error {
	s := recs.Size(r.Len())

	if err := godbase.WriteVal(&s, w); err != nil {
		return err
	}

	for i := r.Iter(); i.Valid(); i=i.Next() {
		if err := cols.Write(i.Key().(cols.Any), i.Val(), w); err != nil {
			return err
		}
	}

	return nil
}
