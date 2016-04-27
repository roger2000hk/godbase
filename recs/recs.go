package recs

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/maps"
)

type Any interface {
	Delete(cols.Any) bool
	Find(cols.Any) (interface{}, bool)
	Get(cols.Any) interface{}
	Id() Id
	Int64(*cols.Int64) int64
	Iter() Iter
	Len() int
	Set(cols.Any, interface{}) interface{}
	SetInt64(*cols.Int64, int64) int64
	SetString(*cols.String, string) string
	SetUId(*cols.UId, godbase.UId) godbase.UId
	String(*cols.String) string
	UId(*cols.UId) godbase.UId
}

type Basic maps.Skip
type Id godbase.UId
type Iter maps.Iter
type Size uint32
type Alloc maps.SkipAlloc

var idCol = cols.NewUId("godbase/id")

func IdCol() *cols.UId {
	return idCol
}

func New(alloc *Alloc) Any {
	return new(Basic).Init(alloc)
}

func NewAlloc(s int) *Alloc {
	return (*Alloc)(maps.NewSkipAlloc(s))
}

func NewId() Id {
	return Id(godbase.NewUId())
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

func (r *Basic) Id() Id {
	return Id(r.UId(idCol))
}

func (r *Basic) Init(alloc *Alloc) *Basic {
	r.asMap().Init((*maps.SkipAlloc)(alloc), 1)
	r.SetUId(idCol, godbase.NewUId())
	return r
}

func (r *Basic) Iter() Iter {
	return Iter(r.asMap().First())
}

func (r *Basic) Int64(c *cols.Int64) int64 {
	return r.Get(c).(int64)
}

func (r *Basic) Len() int {
	return int(r.asMap().Len())
}

func (r *Basic) Set(c cols.Any, v interface{}) interface{} {
	return r.asMap().Set(c, v)
}

func (r *Basic) SetInt64(c *cols.Int64, v int64) int64 {
	return r.Set(c, v).(int64)
}

func (r *Basic) SetString(c *cols.String, v string) string {
	return r.Set(c, v).(string)
}

func (r *Basic) SetUId(c *cols.UId, v godbase.UId) godbase.UId {
	return r.Set(c, v).(godbase.UId)
}

func (r *Basic) String(c *cols.String) string {
	return r.Get(c).(string)
}

func (r *Basic) UId(c *cols.UId) godbase.UId {
	return r.Get(c).(godbase.UId)
}

func (r *Basic) asMap() *maps.Skip {
	return (*maps.Skip)(r)
}
