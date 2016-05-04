package cols

import (
	"fmt"
	"github.com/fncodr/godbase"
)

type RefCol struct {
	UIdCol
}

type RefType struct {
	UIdType
	tbl godbase.Tbl
}

func NewRef(n string, t godbase.Tbl) *RefCol {
	return new(RefCol).Init(n, t)
}

func (c *RefCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *RefType) Encode(_v interface{}) interface{} {
	if v, ok := _v.(godbase.Rec); ok {
		return godbase.UId(v.Id())
	}

	return _v.(godbase.UId)
}

func (c *RefCol) Init(n string, tbl godbase.Tbl) *RefCol {
	c.Basic.Init(n, Ref(tbl))
	return c
}

func (t *RefType) Init(tbl godbase.Tbl) *RefType {
	t.BasicType.Init(fmt.Sprintf("Ref/%v", tbl.Name()))
	t.tbl = tbl
	typeRegistry[t.Name()] = t
	return t
}

func (c *RefCol) Tbl() godbase.Tbl {
	return c.colType.(*RefType).tbl
}
