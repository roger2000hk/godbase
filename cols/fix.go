package cols

import (
	"fmt"
	"github.com/fncodr/godbase"
	"github.com/fncodr/godbase/fix"
	"hash"
	"io"
)

type FixCol struct {
	Basic
}

type FixType struct {
	BasicType
	denom int64
}

func NewFix(n string, m int64) *FixCol {
	return new(FixCol).Init(n, m)
}

func (c *FixCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (t *FixType) AsKey(_ godbase.Rec, _v interface{}) godbase.Key {
	if v, ok := _v.(fix.Val); ok {
		v.Scale(t.denom)
		return godbase.FixKey(v)
	}

	return godbase.Int64Key(_v.(int64))
}

func (c *FixCol) Denom() int64 {
	return c.colType.(*FixType).denom
}

func (t *FixType) Decode(v interface{}) interface{} {
	var res fix.Val
	res.Init(v.(int64), t.denom)
	return res
}

func (t *FixType) Encode(_v interface{}) interface{} {
	if v, ok := _v.(fix.Val); ok {
		return v.Scale(t.denom).Num()
	}

	return _v
}

func (_ *FixType) Eq(_l, _r interface{}) bool {
	return _l.(int64) == _r.(int64)
}

func (_ *FixType) Hash(_ godbase.Rec, _v interface{}, h hash.Hash64) {
	v := _v.(godbase.Int64Key)
	godbase.Write(&v, h)
}

func (c *FixCol) Init(n string, d int64) *FixCol {
	c.Basic.Init(n, new(FixType).Init(d))
	return c
}

func (t *FixType) Init(d int64) *FixType {
	t.BasicType.Init(fmt.Sprintf("Fix(%v)", d))
	t.denom = d
	typeRegistry[t.Name()] = t
	return t
}

func (_ *FixType) Read(_ godbase.Rec, s godbase.ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *FixType) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}
