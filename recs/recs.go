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
	New() Any
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
	SetTime(*cols.Time, time.Time) time.Time
	SetUId(*cols.UId, godbase.UId) godbase.UId
	String(*cols.String) string
	Time(*cols.Time) time.Time
	UId(*cols.UId) godbase.UId
}

type Basic maps.Skip
type Id godbase.UId
type Iter maps.Iter
type Size uint32
type Alloc *maps.SkipAlloc

var createdAtCol = cols.NewTime("godbase/createdAt")
var idCol = cols.NewUId("godbase/id")

func CreatedAtCol() *cols.Time {
	return createdAtCol
}

type IdHash struct {
	imp hash.Hash64
}

func NewIdHash() *IdHash {
	return new(IdHash).Init()
}

func (h *IdHash) Hash(id Id) uint64 {
	h.imp.Reset()
	h.imp.Write(id[:])
	return h.imp.Sum64()
}

func IdCol() *cols.UId {
	return idCol
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

func (r *Basic) CreatedAt() time.Time {
	return r.Time(createdAtCol)
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

func (r *Basic) Init(alloc Alloc) *Basic {
	r.asMap().Init((*maps.SkipAlloc)(alloc), 1)
	r.SetTime(createdAtCol, time.Now())
	r.SetUId(idCol, godbase.NewUId())
	return r
}

func (h *IdHash) Init() *IdHash {
	h.imp = fnv.New64()
	return h
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

func (id Id) Less(other maps.Key) bool {
	oid := other.(Id)
	return bytes.Compare(id[:], oid[:]) < 0
}

func (r *Basic) New() Any {
	return r.asMap().New()
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

func (r *Basic) SetTime(c *cols.Time, v time.Time) time.Time {
	return r.Set(c, v).(time.Time)
}

func (r *Basic) SetUId(c *cols.UId, v godbase.UId) godbase.UId {
	return r.Set(c, v).(godbase.UId)
}

func (r *Basic) String(c *cols.String) string {
	return r.Get(c).(string)
}

func (r *Basic) Time(c *cols.Time) time.Time {
	return r.Get(c).(time.Time)
}

func (r *Basic) UId(c *cols.UId) godbase.UId {
	return r.Get(c).(godbase.UId)
}

func (r *Basic) asMap() *maps.Skip {
	return (*maps.Skip)(r)
}
