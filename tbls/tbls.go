package tbls

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"io"
	"log"
	"time"
)

type Any interface {
	defs.Any
	Add(cols.Any) cols.Any
	Clear()
	Col(string) cols.Any
	Cols() ColIter
	Dump(io.Writer) error
	Get(recs.Id) (recs.Any, error)
	Len() int64
	Reset(recs.Any) (recs.Any, error)
	Read(recs.Any, io.Reader) (recs.Any, error)
	Revision(r recs.Any) int64
	Slurp(io.Reader) error
	Upsert(recs.Any) (recs.Any, error)
	UpsertedAt(recs.Any) time.Time
	Write(recs.Any, io.Writer) error
}

type Basic struct {
	defs.Basic
	cols maps.Skip
	recIdHash recs.IdHash
	recs maps.Hash
	revision cols.Int64Col
	upsertedAt cols.TimeCol
}

type ColIter maps.Iter
type RecNotFound recs.Id

func New(n string, rsc int, ra *maps.SkipAlloc, rls int) Any {
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
	if _, ok := t.cols.Insert(nil, maps.StringKey(c.Name()), c, false); ok {
		return c
	}

	panic(fmt.Sprintf("dup col: %v!", c.Name()))
}

func (t *Basic) Clear() {
	t.recs.While(func (_ maps.Key, v interface{}) bool {
		v.(recs.Any).Clear()
		return true
	})

	t.recs.Clear()
}

func (t *Basic) Dump(w io.Writer) error {
	var err error

	t.While(func (r recs.Any) bool {
		if err = t.Write(r, w); err != nil {
			return false
		}

		return true
	})

	return err
}

func (e RecNotFound) Error() string {
	return fmt.Sprintf("rec not found: %v", e)
}

func (t *Basic) Get(id recs.Id) (recs.Any, error) {
	rr, ok := t.recs.Get(id)
	if !ok {
		return nil, RecNotFound(id)
	}
	
	return rr.(recs.Any).Clone(), nil
}

func (t *Basic) Init(n string, rsc int, ra *maps.SkipAlloc, rls int) *Basic {
	t.Basic.Init(n)
	t.cols.Init(nil, 1)
	t.recIdHash.Init()

	hashRecId := func(_id maps.Key) uint64 {
		id := _id.(recs.Id)
		return t.recIdHash.Hash(id)
	}

	t.recs.Init(maps.NewSkipSlots(rsc, hashRecId, ra, rls))
	t.Add(cols.CreatedAt())
	t.Add(cols.RecId())
	t.Add(t.revision.Init(fmt.Sprintf("%v/revision", n)))
	t.Add(t.upsertedAt.Init(fmt.Sprintf("%v/upsertedAt", n)))
	return t
}

func (t *Basic) Len() int64 {
	return t.recs.Len()
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

		if i, ok := t.cols.Find(nil, maps.StringKey(n), nil); ok {
			c := i.Val().(cols.Any)
			var v interface{}
			
			if v, err = cols.Read(c, r); err != nil {
				return nil, err
			}

			rec.Set(c, v)
		} else {
			log.Printf("col '%v' missing in tbl '%v'", n, t.Name())
 
			var s cols.ValSize

			if s, err = cols.ReadSize(r); err != nil {
				return nil, err
			}

			bs := make([]byte, s)
			if _, err = io.ReadFull(r, bs); err != nil {
				return nil, err
			}
		}
	}

	return rec, nil
}

func (t *Basic) Reset(rec recs.Any) (recs.Any, error) {
	id := rec.Id()
	rr, ok := t.recs.Get(id)
	if !ok {
		return nil, RecNotFound(id)
	}
	
	for i := rr.(recs.Any).Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(cols.Any)
		rec.Set(c, c.CloneVal(i.Val()))
	}
	
	return rec, nil
}

func (t *Basic) Revision(r recs.Any) (v int64) {
	if v, ok := r.Find(&t.revision); ok {
		return v.(int64)
	}
	
	return -1
}

func (t *Basic) Slurp(r io.Reader) error {
	for true {
		rec, err := t.Read(recs.New(nil), r)

		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
		
		t.recs.Set(rec.Id(), rec)
	}

	return nil
}

func (t *Basic) Upsert(rec recs.Any) (recs.Any, error) {
	id := rec.Id()

	if v, ok := rec.Find(&t.revision); ok {
		rec.Set(&t.revision, v.(int64) + 1)
	} else {
		rec.Set(&t.revision, int64(0))
	}

	rec.Set(&t.upsertedAt, time.Now())
	rr := rec.New()
	
	for i := t.cols.First(); i.Valid(); i = i.Next() {
		c := i.Val().(cols.Any)
		if v, ok := rec.Find(c); ok {
			rr.Set(c, c.CloneVal(v))
		}
	}
	
	t.recs.Set(id, rr)
	return rec, nil
}

func (t *Basic) UpsertedAt(r recs.Any) (res time.Time) {
	if v, ok := r.Find(&t.upsertedAt); ok {
		return v.(time.Time)
	}

	return res	
}

func (t *Basic) While(fn recs.TestFn) bool {
	return t.recs.While(func (_ maps.Key, v interface{}) bool {
		return fn(v.(recs.Any))
	})
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
