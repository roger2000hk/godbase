package tbls

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/defs"
	"github.com/fncodr/godbase/idxs"
	"github.com/fncodr/godbase/maps"
	"github.com/fncodr/godbase/recs"
	"io"
	"log"
	"time"
)

type Basic struct {
	defs.Basic
	cols maps.Sort
	onUpsert godbase.Evt
	recIdHash godbase.UIdHash
	recs maps.Hash
	revision cols.Int64Col
	upsertedAt cols.TimeCol
}

type InsertFn func(godbase.Rec) error
type RecNotFound godbase.UId

func AddBool(t godbase.Tbl, n string) *cols.BoolCol {
	return t.Add(cols.NewBool(n)).(*cols.BoolCol)
}

func AddFix(t godbase.Tbl, n string, d int64) *cols.FixCol {
	return t.Add(cols.NewFix(n, d)).(*cols.FixCol)
}

func AddHashIdx(t godbase.Tbl, n string, cs []godbase.Col, u bool, sc int, a *maps.SlabAlloc, 
	ls int) godbase.Idx {
	i := idxs.NewHash(n, cs, u, sc, a, ls)
	AddIdx(t, i)
	return i
}

func AddIdx(t godbase.Tbl, i godbase.Idx) godbase.Idx {
	OnUpsert(t, i, func (rec godbase.Rec) error { 
		if prev, err := t.Get(rec.Id()); err != nil {
			if _, ok := err.(RecNotFound); !ok {
				return err
			}
		} else {
			i.Delete(prev)
		}

		_, err := i.Insert(rec)
		return err
	})

	return i
}

func AddInt64(t godbase.Tbl, n string) *cols.Int64Col {
	return t.Add(cols.NewInt64(n)).(*cols.Int64Col)
}

func AddRef(t godbase.Tbl, n string, rt godbase.Tbl) *cols.RefCol {
	return t.Add(cols.NewRef(n, rt)).(*cols.RefCol)
}

func AddSortIdx(t godbase.Tbl, n string, cs []godbase.Col, u bool, a *maps.SlabAlloc, 
	ls int) godbase.Idx {
	i := idxs.NewSort(n, cs, u, a, ls)
	AddIdx(t, i)
	return i
}

func AddString(t godbase.Tbl, n string) *cols.StringCol {
	return t.Add(cols.NewString(n)).(*cols.StringCol)
}

func AddTime(t godbase.Tbl, n string) *cols.TimeCol {
	return t.Add(cols.NewTime(n)).(*cols.TimeCol)
}

func AddUId(t godbase.Tbl, n string) *cols.UIdCol {
	return t.Add(cols.NewUId(n)).(*cols.UIdCol)
}

func AddUnion(t godbase.Tbl, n string, fn cols.UnionTypeFn) *cols.UnionCol {
	return t.Add(cols.NewUnion(n, fn)).(*cols.UnionCol)
}

func New(n string, rsc int, ra *maps.SlabAlloc, rls int) godbase.Tbl {
	return new(Basic).Init(n, rsc, ra, rls)
}

func OnUpsert(t godbase.Tbl, sub godbase.EvtSub, fn InsertFn) {
	t.OnUpsert().Subscribe(sub, func(args...interface{}) error {
		return fn(args[0].(godbase.Rec))
	})
}

func (t *Basic) Col(n string) godbase.Col {
	if c, ok := t.cols.Get(godbase.StringKey(n)); ok {
		return c.(godbase.Col)
	}
	
	panic(fmt.Sprintf("col not found: %v", n))
}

func (t *Basic) Cols() godbase.Iter {
	return t.cols.First()
}

func (t *Basic) Add(c godbase.Col) godbase.Col {
	if _, ok := t.cols.Insert(nil, godbase.StringKey(c.Name()), c, false); ok {
		return c
	}

	panic(fmt.Sprintf("dup col: %v!", c.Name()))
}

func (t *Basic) Clear() {
	t.recs.While(func (_ godbase.Key, v interface{}) bool {
		v.(godbase.Rec).Clear()
		return true
	})

	t.recs.Clear()
}

func (t *Basic) Dump(w io.Writer) error {
	var err error

	t.While(func (r godbase.Rec) bool {
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

func (t *Basic) Get(id godbase.UId) (godbase.Rec, error) {
	rr, ok := t.recs.Get(godbase.UIdKey(id))
	if !ok {
		return nil, RecNotFound(id)
	}
	
	return rr.(godbase.Rec).Clone(), nil
}

func (t *Basic) Init(n string, rsc int, ra *maps.SlabAlloc, rls int) *Basic {
	t.Basic.Init(n)
	t.cols.Init(nil, 1)
	t.onUpsert.Init()
	t.recIdHash.Init()
	
	hashRecId := func(_id godbase.Key) uint64 {
		id := godbase.UId(_id.(godbase.UIdKey))
		return t.recIdHash.Hash(id)
	}

	t.recs.Init(maps.NewSlabSlots(rsc, hashRecId, ra, rls))
	t.Add(cols.CreatedAt())
	t.Add(cols.RecId())
	t.Add(t.revision.Init(fmt.Sprintf("%v/revision", n)))
	t.Add(t.upsertedAt.Init(fmt.Sprintf("%v/upsertedAt", n)))
	return t
}

func (t *Basic) Len() int64 {
	return t.recs.Len()
}

func (self *Basic) OnUpsert() *godbase.Evt {
	return &self.onUpsert
}

func (t *Basic) Read(rec godbase.Rec, r io.Reader) (godbase.Rec, error) {
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

		if i, ok := t.cols.Find(nil, godbase.StringKey(n), nil); ok {
			c := i.Val().(godbase.Col)
			var v interface{}
			
			if v, err = cols.Read(rec, c, r); err != nil {
				return nil, err
			}

			rec.Set(c, v)
		} else {
			log.Printf("col '%v' missing in tbl '%v'", n, t.Name())
 
			var s godbase.ValSize

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

func (t *Basic) Reset(rec godbase.Rec) (godbase.Rec, error) {
	id := rec.Id()
	rr, ok := t.recs.Get(godbase.UIdKey(id))
	if !ok {
		return nil, RecNotFound(id)
	}
	
	for i := rr.(godbase.Rec).Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(godbase.Col)
		rec.Set(c, c.CloneVal(i.Val()))
	}
	
	return rec, nil
}

func (t *Basic) Revision(r godbase.Rec) (v int64) {
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
		
		t.recs.Set(godbase.UIdKey(rec.Id()), rec)
	}

	return nil
}

func (t *Basic) Upsert(rec godbase.Rec) (godbase.Rec, error) {
	id := rec.Id()

	if v, ok := rec.Find(&t.revision); ok {
		rec.Set(&t.revision, v.(int64) + 1)
	} else {
		rec.Set(&t.revision, int64(0))
	}

	rec.Set(&t.upsertedAt, time.Now())
	rr := rec.New()
	
	for i := t.cols.First(); i.Valid(); i = i.Next() {
		c := i.Val().(godbase.Col)
		if v, ok := rec.Find(c); ok {
			rr.Set(c, c.CloneVal(v))
		}
	}
	
	if err := t.onUpsert.Publish(rr); err != nil {
		return nil, err
	}

	t.recs.Set(godbase.UIdKey(id), rr)
	return rec, nil
}

func (t *Basic) UpsertedAt(r godbase.Rec) (res time.Time) {
	if v, ok := r.Find(&t.upsertedAt); ok {
		return v.(time.Time)
	}

	return res	
}

func (t *Basic) While(fn recs.TestFn) bool {
	return t.recs.While(func (_ godbase.Key, v interface{}) bool {
		return fn(v.(godbase.Rec))
	})
}

func (t *Basic) Write(rec godbase.Rec, w io.Writer) error {
	s := recs.Size(rec.Len())

	if err := godbase.Write(&s, w); err != nil {
		return err
	}

	for i := rec.Iter(); i.Valid(); i=i.Next() {
		if err := cols.Write(rec, i.Key().(godbase.Col), i.Val(), w); err != nil {
			return err
		}
	}

	return nil
}
