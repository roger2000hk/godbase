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

type ColIter maps.Iter

type Any interface {
	defs.Any
	Add(cols.Any) cols.Any
	Col(n string) cols.Any
	Cols() ColIter
	Read(rec recs.Any, r io.Reader) (recs.Any, error)
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
