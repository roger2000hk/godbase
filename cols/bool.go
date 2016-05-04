package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
)

type BoolCol struct {
	Basic
}

type BoolType struct {
	BasicType
}

func NewBool(n string) *BoolCol {
	return new(BoolCol).Init(n)
}

func (c *BoolCol) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *BoolType) AsKey(_ godbase.Rec, v interface{}) godbase.Key {
	return godbase.BoolKey(v.(bool))
}

func (_ *BoolType) Hash(_ godbase.Rec, _v interface{}, h hash.Hash64) {
	v := _v.(godbase.BoolKey)
	godbase.Write(&v, h)
}

func (c *BoolCol) Init(n string) *BoolCol {
	c.Basic.Init(n, Bool())
	return c
}

func (t *BoolType) Init(n string) *BoolType {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	return t
}

func (_ *BoolType) Read(_ godbase.Rec, _ godbase.ValSize, r io.Reader) (interface{}, error) {
	var v byte

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v == 1, nil
}

func (_ *BoolType) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	v := byte(0)
	if _v.(bool) {
		v = 1
	} 

	return WriteBinVal(1, &v, w)
}
