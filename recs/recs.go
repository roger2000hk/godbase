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
	SetInt64(*cols.Int64, int64) int64
	SetString(*cols.String, string) string
	SetUId(*cols.UId, godbase.UId) godbase.UId
	String(*cols.String) string
	UId(*cols.UId) godbase.UId
}

type Id godbase.UId

type Basic maps.Skip

var idCol = cols.NewUId("godbase/id")

func IdCol() *cols.UId {
	return idCol
}

func New() Any {
	return new(Basic).Init()
}

func NewId() Id {
	return Id(godbase.NewUId())
}

func (r *Basic) AsMap() *maps.Skip {
	return (*maps.Skip)(r)
}

func (r *Basic) Delete(c cols.Any) bool {
	_, cnt := r.AsMap().Delete(nil, nil, c, nil)
	return cnt == 1
}

func (r *Basic) Find(c cols.Any) (interface{}, bool) {
	if v, ok := r.AsMap().Get(c); ok {
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
	if v, ok := r.Find(IdCol()); ok {
		return v.(Id)
	}

	res := NewId()
	r.SetUId(IdCol(), godbase.UId(res))
	return res
}

func (r *Basic) Init() *Basic {
	r.AsMap().Init(nil, 1)
	return r
}

func (r *Basic) Int64(c *cols.Int64) int64 {
	return r.Get(c).(int64)
}

func (r *Basic) SetInt64(c *cols.Int64, v int64) int64 {
	return r.AsMap().Set(c, v).(int64)
}

func (r *Basic) SetString(c *cols.String, v string) string {
	return r.AsMap().Set(c, v).(string)
}

func (r *Basic) SetUId(c *cols.UId, v godbase.UId) godbase.UId {
	return r.AsMap().Set(c, v).(godbase.UId)
}

func (r *Basic) String(c *cols.String) string {
	return r.Get(c).(string)
}

func (r *Basic) UId(c *cols.UId) godbase.UId {
	return r.Get(c).(godbase.UId)
}
