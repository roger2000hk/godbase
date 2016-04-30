package recs

import (
	"bytes"
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/maps"
	"hash"
	"hash/fnv"
	"time"
)

type Any interface {
	Bool(*cols.BoolCol) bool
	Clear()
	Clone() Any
	Fix(*cols.FixCol) fix.Val
	Delete(cols.Any) bool
	Eq(Any) bool
	Find(cols.Any) (interface{}, bool)
	Get(cols.Any) interface{}
	Id() Id
	Int64(*cols.Int64Col) int64
	Iter() maps.Iter
	Len() int
	New() Any
	Set(cols.Any, interface{}) Any
	SetBool(*cols.BoolCol, bool) Any
	SetFix(*cols.FixCol, fix.Val) Any
	SetInt64(*cols.Int64Col, int64) Any
	SetString(*cols.StringCol, string) Any
	SetTime(*cols.TimeCol, time.Time) Any
	SetUId(*cols.UIdCol, godbase.UId) Any
	String(*cols.StringCol) string
	Time(*cols.TimeCol) time.Time
	UId(*cols.UIdCol) godbase.UId
}

type IdHash struct {
	imp hash.Hash64
}

type Basic maps.Sort
type Id godbase.UId
type NotFound Id
type Size uint32
type TestFn func(Any) bool

func BasicNew(a *maps.SlabAlloc) Any {
	return new(Basic).BasicInit(a)
}

func Init(id Id, a *maps.SlabAlloc) Any {
	r := new(Basic).BasicInit(a)
	r.SetUId(cols.RecId(), godbase.UId(id))
	return r
}

func New(a *maps.SlabAlloc) Any {
	return new(Basic).Init(a)
}

func NewId() Id {
	return Id(godbase.NewUId())
}

func NewIdHash() *IdHash {
	return new(IdHash).Init()
}

func (e NotFound) Error() string {
	return fmt.Sprintf("rec not found: %v", e)
}

func (r *Basic) Bool(c *cols.BoolCol) bool {
	return r.Get(c).(bool)
}

func (r *Basic) BasicInit(a *maps.SlabAlloc) *Basic {
	r.asMap().Init(a, 1)
	return r
}

func (r *Basic) Clear() {
	r.asMap().Clear()
}

func (r *Basic) CreatedAt() time.Time {
	return r.Time(cols.CreatedAt())
}

func (r *Basic) Clone() Any {
	res := r.New()

	for i := r.Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(cols.Any)
		res.Set(c, c.CloneVal(i.Val()))
	}
	
	return res
}

func (r *Basic) Fix(c *cols.FixCol) (res fix.Val) {
	return r.Get(c).(fix.Val)
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
		return c.Decode(v)
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

func (r *Basic) Init(a *maps.SlabAlloc) *Basic {
	r.BasicInit(a)
	r.SetTime(cols.CreatedAt(), time.Now())
	r.SetUId(cols.RecId(), godbase.NewUId())
	return r
}

func (h *IdHash) Init() *IdHash {
	h.imp = fnv.New64()
	return h
}

func (r *Basic) Iter() maps.Iter {
	return r.asMap().First()
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
	return (*Basic)(r.asMap().New().(*maps.Sort))
}

func (r *Basic) Set(c cols.Any, v interface{}) Any {
	r.asMap().Set(c, c.Encode(v))
	return r
}

func (r *Basic) SetBool(c *cols.BoolCol, v bool) Any {
	return r.Set(c, v)
}

func (r *Basic) SetFix(c *cols.FixCol, v fix.Val) Any {
	return r.Set(c, v)
}

func (r *Basic) SetInt64(c *cols.Int64Col, v int64) Any {
	return r.Set(c, v)
}

func (r *Basic) SetString(c *cols.StringCol, v string) Any {
	return r.Set(c, v)
}

func (r *Basic) SetTime(c *cols.TimeCol, v time.Time) Any {
	return r.Set(c, v)
}

func (r *Basic) SetUId(c *cols.UIdCol, v godbase.UId) Any {
	return r.Set(c, v)
}

func (r *Basic) String(c *cols.StringCol) string {
	return r.Get(c).(string)
}

func (id Id) String() string {
	return godbase.UId(id).String()
}

func (r *Basic) Time(c *cols.TimeCol) time.Time {
	return r.Get(c).(time.Time)
}

func (r *Basic) UId(c *cols.UIdCol) godbase.UId {
	return r.Get(c).(godbase.UId)
}

func (r *Basic) asMap() *maps.Sort {
	return (*maps.Sort)(r)
}
