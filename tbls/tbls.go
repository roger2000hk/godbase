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
	mapAlloc *maps.SlabAlloc
	onDelete, onDrop, onLoad, onUpsert godbase.Evt
	recIdHash godbase.UIdHash
	recs maps.Hash
	revision cols.Int64Col
	upsertedAt cols.TimeCol
}

type OnDeleteFn func(godbase.Cx, godbase.Rec) error
type OnDropFn func(godbase.Cx, godbase.Rec) error
type OnLoadFn func(godbase.Cx, godbase.Rec) error
type OnUpsertFn func(godbase.Cx, godbase.Rec) error

func AddBool(t godbase.Tbl, n string) *cols.BoolCol {
	return t.AddCol(cols.NewBool(n)).(*cols.BoolCol)
}

func AddFix(t godbase.Tbl, n string, d int64) *cols.FixCol {
	return t.AddCol(cols.NewFix(n, d)).(*cols.FixCol)
}

func AddHashIdx(t godbase.Tbl, n string, cs []godbase.Col, u bool, sc int, a *maps.SlabAlloc, 
	ls int) godbase.Idx {
	i := idxs.NewHash(n, cs, u, sc, a, ls)
	t.AddIdx(i)
	return i
}

func AddInt64(t godbase.Tbl, n string) *cols.Int64Col {
	return t.AddCol(cols.NewInt64(n)).(*cols.Int64Col)
}

func AddRef(t godbase.Tbl, n string, rt godbase.Tbl) *cols.RefCol {
	return t.AddCol(cols.NewRef(n, rt)).(*cols.RefCol)
}

func AddSortIdx(t godbase.Tbl, n string, cs []godbase.Col, u bool, a *maps.SlabAlloc, 
	ls int) godbase.Idx {
	i := idxs.NewSort(n, cs, u, a, ls)
	t.AddIdx(i)
	return i
}

func AddStr(t godbase.Tbl, n string) *cols.StrCol {
	return t.AddCol(cols.NewStr(n)).(*cols.StrCol)
}

func AddSuffixIdx(t godbase.Tbl, n string, cs []*cols.StrCol, u bool, a *maps.SlabAlloc, 
	ls int) godbase.Idx {
	i := idxs.NewSuffix(n, cs, u, a, ls)
	t.AddIdx(i)
	return i
}

func AddTime(t godbase.Tbl, n string) *cols.TimeCol {
	return t.AddCol(cols.NewTime(n)).(*cols.TimeCol)
}

func AddUId(t godbase.Tbl, n string) *cols.UIdCol {
	return t.AddCol(cols.NewUId(n)).(*cols.UIdCol)
}

func AddUnion(t godbase.Tbl, n string, fn cols.UnionTypeFn) *cols.UnionCol {
	return t.AddCol(cols.NewUnion(n, fn)).(*cols.UnionCol)
}

func New(n string, ds []godbase.TblDef, rsc int, ma *maps.SlabAlloc, rls int) godbase.Tbl {
	return new(Basic).Init(n, ds, rsc, ma, rls)
}

func OnDelete(t godbase.Tbl, sub godbase.EvtSub, fn OnDeleteFn) {
	t.OnDelete().Subscribe(sub, func(args...interface{}) error {
		return fn(args[0].(godbase.Cx), args[1].(godbase.Rec))
	})
}

func OnDrop(t godbase.Tbl, sub godbase.EvtSub, fn OnDropFn) {
	t.OnDrop().Subscribe(sub, func(args...interface{}) error {
		return fn(args[0].(godbase.Cx), args[1].(godbase.Rec))
	})
}

func OnLoad(t godbase.Tbl, sub godbase.EvtSub, fn OnLoadFn) {
	t.OnLoad().Subscribe(sub, func(args...interface{}) error {
		return fn(args[0].(godbase.Cx), args[1].(godbase.Rec))
	})
}

func OnUpsert(t godbase.Tbl, sub godbase.EvtSub, fn OnUpsertFn) {
	t.OnUpsert().Subscribe(sub, func(args...interface{}) error {
		return fn(args[0].(godbase.Cx), args[1].(godbase.Rec))
	})
}

func (t *Basic) AddIdx(i godbase.Idx) godbase.Idx {
	OnDelete(t, i, func (cx godbase.Cx, rec godbase.Rec) error {
		return i.Delete(nil, rec)
	})

	OnDrop(t, i, func (cx godbase.Cx, rec godbase.Rec) error {
		return i.Drop(nil, rec)
	})

	OnUpsert(t, i, func (cx godbase.Cx, rec godbase.Rec) error {
		if prev, err := t.Reset(recs.New(rec.Id())); err != nil {
			if _, ok := err.(recs.NotFound); !ok {
				return err
			}
		} else {
			i.Delete(nil, prev)
		}

		_, err := i.Insert(nil, rec)
		return err
	})

	return i
}

func (t *Basic) Col(n string) godbase.Col {
	if c, ok := t.cols.Get(godbase.StrKey(n)); ok {
		return c.(godbase.Col)
	}
	
	panic(fmt.Sprintf("col not found: %v", n))
}

func (t *Basic) Cols() godbase.Iter {
	return t.cols.First()
}

func (t *Basic) AddCol(c godbase.Col) godbase.Col {
	if _, ok := t.cols.Insert(nil, godbase.StrKey(c.Name()), c, false); ok {
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

func (t *Basic) Delete(cx godbase.Cx, id godbase.UId) error {
	k := godbase.UIdKey(id)
	i, ok := t.recs.Find(nil, k, nil)

	if !ok {
		return recs.NotFound(id)
	}

	if err := t.onDelete.Publish(cx, i.Val().(godbase.Rec)); err != nil {
		return err
	}

	if _, ok := t.recs.Delete(i, nil, godbase.UIdKey(id), nil); ok != 1{
		panic(fmt.Sprintf("delete failed: %v", id))
	}

	return nil
}

func (t *Basic) Drop(cx godbase.Cx, id godbase.UId) error {
	k := godbase.UIdKey(id)
	i, ok := t.recs.Find(nil, k, nil)

	if !ok {
		return recs.NotFound(id)
	}

	if err := t.onDrop.Publish(cx, i.Val().(godbase.Rec)); err != nil {
		return err
	}

	if _, ok := t.recs.Delete(i, nil, godbase.UIdKey(id), nil); ok != 1 {
		panic(fmt.Sprintf("delete failed: %v", id))
	}

	return nil
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

func (self *Basic) Exists(id godbase.UId) bool {
	_, ok := self.recs.Find(nil, godbase.UIdKey(id), nil)
	return ok
}

func (t *Basic) Init(n string, ds []godbase.TblDef, rsc int, ma *maps.SlabAlloc, rls int) *Basic {
	t.Basic.Init(n)
	t.cols.Init(nil, 1)
	t.mapAlloc = ma
	t.onDelete.Init()
	t.onDrop.Init()
	t.onLoad.Init()
	t.onUpsert.Init()
	t.recIdHash.Init()
	
	hashRecId := func(_id godbase.Key) uint64 {
		id := godbase.UId(_id.(godbase.UIdKey))
		return t.recIdHash.Hash(id)
	}

	t.recs.Init(maps.NewSlabSlots(rsc, hashRecId, ma, rls))
	t.AddCol(cols.CreatedAt())
	t.AddCol(cols.RecId())
	t.AddCol(t.revision.Init(fmt.Sprintf("%v/revision", n)))
	t.AddCol(t.upsertedAt.Init(fmt.Sprintf("%v/upsertedAt", n)))

	if ds != nil {
		for _, d := range ds {
			d.AddToTbl(t)
		}
	}

	return t
}

func (t *Basic) Len() int64 {
	return t.recs.Len()
}

func (self *Basic) OnDelete() *godbase.Evt {
	return &self.onDelete
}

func (self *Basic) OnDrop() *godbase.Evt {
	return &self.onDrop
}

func (self *Basic) OnLoad() *godbase.Evt {
	return &self.onLoad
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

		if i, ok := t.cols.Find(nil, godbase.StrKey(n), nil); ok {
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
		return nil, recs.NotFound(id)
	}
	
	rr.(godbase.Rec).While(func(c godbase.Col, v interface{}) bool {
		rec.Set(c, c.CloneVal(v))
		return true
	})
	
	return rec, nil
}

func (t *Basic) Revision(r godbase.Rec) (v int64) {
	if v, ok := r.Find(&t.revision); ok {
		return v.(int64)
	}
	
	return -1
}

func (t *Basic) Load(cx godbase.Cx, rec godbase.Rec) (godbase.Rec, error) {
	if err := t.onLoad.Publish(cx, rec); err != nil {
		return nil, err
	}
	
	t.recs.Set(godbase.UIdKey(rec.Id()), rec)
	return rec, nil
}

func (t *Basic) Slurp(cx godbase.Cx, r io.Reader) error {
	for true {
		rec, err := t.Read(new(recs.Basic), r)

		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}
		
		t.Load(cx, rec)
	}

	return nil
}

func (t *Basic) Upsert(cx godbase.Cx, rec godbase.Rec) (godbase.Rec, error) {
	id := rec.Id()

	if v, ok := rec.Find(&t.revision); ok {
		rec.Set(&t.revision, v.(int64) + 1)
	} else {
		rec.Set(&t.revision, int64(0))
	}

	rec.Set(&t.upsertedAt, time.Now())
	rr := new(recs.Basic)
	
	for i := t.cols.First(); i.Valid(); i = i.Next() {
		c := i.Val().(godbase.Col)
		if v, ok := rec.Find(c); ok {
			rr.Set(c, c.CloneVal(v))
		}
	}
	
	if err := t.onUpsert.Publish(cx, rr); err != nil {
		return nil, err
	}

	if i, ok := t.recs.Find(nil, godbase.UIdKey(id), nil); ok {
		i.Val().(godbase.Rec).Clear()
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

	var err error

	rec.While(func(c godbase.Col, v interface{}) bool {
		if err = cols.Write(rec, c, v, w); err != nil {
			return false
		}
		
		return true
	})

	return err
}
