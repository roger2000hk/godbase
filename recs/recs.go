package recs

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
	"hash"
	"hash/fnv"
	"time"
)

type Any interface {
	Eq(Any) bool
	New() Any
	Delete(cols.Any) bool
	Find(cols.Any) (interface{}, bool)
	Get(cols.Any) interface{}
	Id() Id
	Int64(*cols.Int64Col) int64
	Iter() Iter
	Len() int
	Set(cols.Any, interface{}) interface{}
	SetInt64(*cols.Int64Col, int64) int64
	SetString(*cols.StringCol, string) string
	SetTime(*cols.TimeCol, time.Time) time.Time
	SetUId(*cols.UIdCol, godbase.UId) godbase.UId
	String(*cols.StringCol) string
	Time(*cols.TimeCol) time.Time
	UId(*cols.UIdCol) godbase.UId
}

type IdHash struct {
	imp hash.Hash64
}

type Alloc *maps.SkipAlloc
type Basic maps.Skip
type Id godbase.UId
type Iter maps.Iter
type Size uint32
type TestFn func(Any) bool

func BasicNew(alloc Alloc) Any {
	return new(Basic).BasicInit(alloc)
}

func Init(id Id, alloc Alloc) Any {
	r := new(Basic).BasicInit(alloc)
	r.SetUId(cols.RecId(), godbase.UId(id))
	return r
}

func New(alloc Alloc) Any {
	return new(Basic).Init(alloc)
}

func NewAlloc(s int) Alloc {
	return Alloc(maps.NewSkipAlloc(s))
}

func NewId() Id {
	return Id(godbase.NewUId())
}

func NewIdHash() *IdHash {
	return new(IdHash).Init()
}

func (r *Basic) BasicInit(alloc Alloc) *Basic {
	r.asMap().Init((*maps.SkipAlloc)(alloc), 1)
	return r
}

func (r *Basic) CreatedAt() time.Time {
	return r.Time(cols.CreatedAt())
}

func (r *Basic) Delete(c cols.Any) bool {
	_, cnt := r.asMap().Delete(nil, nil, c, nil)
	return cnt == 1
}

func (r *Basic) Find(c cols.Any) (interface{}, bool) {
	if v, ok := r.asMap().Get(c); ok {
		return v, true
	}
	
	return nil, false
}

func (r *Basic) Get(c cols.Any) interface{} {
	if v, ok := r.Find(c); ok {
		return v
	}

	panic(fmt.Sprintf("field not found: %v", c.Name()))
}

func (r *Basic) Eq(other Any) bool {
	for i := r.Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(cols.Any)
		if ov, ok := other.Find(c); !ok || !c.Eq(ov, i.Val()) {
			return false
		}
	}

	return true
}

func (h *IdHash) Hash(id Id) uint64 {
	h.imp.Reset()
	h.imp.Write(id[:])
	return h.imp.Sum64()
}

func (r *Basic) Id() Id {
	return Id(r.UId(cols.RecId()))
}

func (r *Basic) Init(alloc Alloc) *Basic {
	r.BasicInit(alloc)
	r.SetTime(cols.CreatedAt(), time.Now())
	r.SetUId(cols.RecId(), godbase.NewUId())
	return r
}

func (h *IdHash) Init() *IdHash {
	h.imp = fnv.New64()
	return h
}

func (r *Basic) Iter() Iter {
	return Iter(r.asMap().First())
}

func (r *Basic) Int64(c *cols.Int64Col) int64 {
	return r.Get(c).(int64)
}

func (r *Basic) Len() int {
	return int(r.asMap().Len())
}

func (id Id) Less(other maps.Key) bool {
	oid := other.(Id)
	return bytes.Compare(id[:], oid[:]) < 0
}

func (r *Basic) New() Any {
	return (*Basic)(r.asMap().New().(*maps.Skip))
}

func (r *Basic) Set(c cols.Any, v interface{}) interface{} {
	r.asMap().Set(c, v)
	return v
}

func (r *Basic) SetInt64(c *cols.Int64Col, v int64) int64 {
	return r.Set(c, v).(int64)
}

func (r *Basic) SetString(c *cols.StringCol, v string) string {
	return r.Set(c, v).(string)
}

func (r *Basic) SetTime(c *cols.TimeCol, v time.Time) time.Time {
	return r.Set(c, v).(time.Time)
}

func (r *Basic) SetUId(c *cols.UIdCol, v godbase.UId) godbase.UId {
	return r.Set(c, v).(godbase.UId)
}

func (r *Basic) String(c *cols.StringCol) string {
	return r.Get(c).(string)
}

func (r *Basic) Time(c *cols.TimeCol) time.Time {
	return r.Get(c).(time.Time)
}

func (r *Basic) UId(c *cols.UIdCol) godbase.UId {
	return r.Get(c).(godbase.UId)
}

func (r *Basic) asMap() *maps.Skip {
	return (*maps.Skip)(r)
}
