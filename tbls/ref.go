package tbls

import (
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/cols"
	"github.com/fncodr/godbase/recs"
)

type RefType struct {
	cols.UIdType
	tbl Any
}

type RefCol struct {
	cols.UIdCol
}

func NewRef(n string, t Any) *RefCol {
	return new(RefCol).Init(n, t)
}

func Ref(tbl Any) cols.Type {
	return new(RefType).Init(tbl)
}

func (rt *RefType) Decode(_v interface{}) interface{} {
	if v, err := rt.tbl.Get(recs.Id(_v.(godbase.UId))); err != nil {
		return v
	}

	return nil
}

func (rt *RefType) Encode(v interface{}) interface{} {
	return godbase.UId(v.(recs.Any).Id())
}

func (rt *RefType) Init(t Any) *RefType {
	rt.tbl = t
	return rt
}

func (c *RefCol) Init(n string, t Any) *RefCol {
	c.Basic.Init(n, Ref(t))
	return c
}



