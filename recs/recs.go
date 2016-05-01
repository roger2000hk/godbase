package recs

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/fix"
	"github.com/fncodr/godbase/maps"
	"time"
)

type Basic maps.Sort
type NotFound godbase.UId
type Size uint32
type TestFn func(godbase.Rec) bool

func BasicNew(a *maps.SlabAlloc) godbase.Rec {
	return new(Basic).BasicInit(a)
}

func Init(id godbase.UId, a *maps.SlabAlloc) godbase.Rec {
	r := new(Basic).BasicInit(a)
	SetUId(r, cols.RecId(), id)
	return r
}

func New(a *maps.SlabAlloc) godbase.Rec {
	return new(Basic).Init(a)
}

func (e NotFound) Error() string {
	return fmt.Sprintf("rec not found: %v", e)
}

func Bool(r godbase.Rec, c *cols.BoolCol) bool {
	return r.Get(c).(bool)
}

func Fix(r godbase.Rec, c *cols.FixCol) (res fix.Val) {
	return r.Get(c).(fix.Val)
}

func Int64(r godbase.Rec, c *cols.Int64Col) int64 {
	return r.Get(c).(int64)
}

func Ref(r godbase.Rec, c *cols.RefCol, res godbase.Rec) (godbase.Rec, error) {
	res.Set(cols.RecId(), r.Get(c).(godbase.UId)) 
	return c.Tbl().Reset(res);
}

func SetBool(r godbase.Rec, c *cols.BoolCol, v bool) bool {
	return r.Set(c, v).(bool)
}

func SetFix(r godbase.Rec, c *cols.FixCol, v fix.Val) fix.Val {
	return r.Set(c, v).(fix.Val)
}

func SetInt64(r godbase.Rec, c *cols.Int64Col, v int64) int64 {
	return r.Set(c, v).(int64)
}

func SetRef(r godbase.Rec, c *cols.RefCol, v godbase.Rec) godbase.Rec {
	return r.Set(c, v).(godbase.Rec)
}

func SetString(r godbase.Rec, c *cols.StringCol, v string) string {
	return r.Set(c, v).(string)
}

func SetTime(r godbase.Rec, c *cols.TimeCol, v time.Time) time.Time {
	return r.Set(c, v).(time.Time)
}

func SetUId(r godbase.Rec, c *cols.UIdCol, v godbase.UId) godbase.UId {
	return r.Set(c, v).(godbase.UId)
}

func String(r godbase.Rec, c *cols.StringCol) string {
	return r.Get(c).(string)
}

func Time(r godbase.Rec, c *cols.TimeCol) time.Time {
	return r.Get(c).(time.Time)
}

func UId(r godbase.Rec, c *cols.UIdCol) godbase.UId {
	return r.Get(c).(godbase.UId)
}

func (r *Basic) BasicInit(a *maps.SlabAlloc) *Basic {
	r.asMap().Init(a, 1)
	return r
}

func (r *Basic) Clear() {
	r.asMap().Clear()
}

func (r *Basic) CreatedAt() time.Time {
	return Time(r, cols.CreatedAt())
}

func (r *Basic) Clone() godbase.Rec {
	res := r.New()

	for i := r.Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(godbase.Col)
		res.Set(c, c.CloneVal(i.Val()))
	}
	
	return res
}

func (r *Basic) Delete(c godbase.Col) bool {
	_, cnt := r.asMap().Delete(nil, nil, c, nil)
	return cnt == 1
}

func (r *Basic) Find(c godbase.Col) (interface{}, bool) {
	if v, ok := r.asMap().Get(c); ok {
		return v, true
	}
	
	return nil, false
}

func (r *Basic) Get(c godbase.Col) interface{} {
	if v, ok := r.Find(c); ok {
		return c.Decode(v)
	}

	panic(fmt.Sprintf("field not found: %v", c.Name()))
}

func (r *Basic) Eq(other godbase.Rec) bool {
	for i := r.Iter(); i.Valid(); i = i.Next() {
		c := i.Key().(godbase.Col)
		if ov, ok := other.Find(c); !ok || !c.Eq(ov, i.Val()) {
			return false
		}
	}

	return true
}

func (r *Basic) Id() godbase.UId {
	return UId(r, cols.RecId())
}

func (r *Basic) Init(a *maps.SlabAlloc) *Basic {
	r.BasicInit(a)
	SetTime(r, cols.CreatedAt(), time.Now())
	SetUId(r, cols.RecId(), godbase.NewUId())
	return r
}

func (r *Basic) Iter() godbase.Iter {
	return r.asMap().First()
}

func (r *Basic) Len() int {
	return int(r.asMap().Len())
}

func (r *Basic) New() godbase.Rec {
	return (*Basic)(r.asMap().New().(*maps.Sort))
}

func (r *Basic) Set(c godbase.Col, v interface{}) interface{} {
	r.asMap().Set(c, c.Encode(v))
	return v
}

func (r *Basic) asMap() *maps.Sort {
	return (*maps.Sort)(r)
}
