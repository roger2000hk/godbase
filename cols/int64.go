package cols

import (
	"github.com/fncodr/godbase"
	"hash"
	"io"
)

type Int64Col struct {
	Basic
}

type Int64Type struct {
	BasicType
}

func NewInt64(n string) *Int64Col {
	return new(Int64Col).Init(n)
}

func (c *Int64Col) AddToTbl(t godbase.Tbl) {
	t.AddCol(c)
}

func (_ *Int64Type) AsKey(_ godbase.Rec, v interface{}) godbase.Key {
	return godbase.Int64Key(v.(int64))
}

func (_ *Int64Type) Hash(_ godbase.Rec, _v interface{}, h hash.Hash64) {
	v := _v.(godbase.Int64Key)
	godbase.Write(&v, h)
}

func (c *Int64Col) Init(n string) *Int64Col {
	c.Basic.Init(n, Int64())
	return c
}

func (t *Int64Type) Init(n string) *Int64Type {
	t.BasicType.Init(n)
	typeRegistry[n] = t
	return t
}

func (_ *Int64Type) Read(_ godbase.Rec, _ godbase.ValSize, r io.Reader) (interface{}, error) {
	var v int64

	if err := godbase.Read(&v, r); err != nil {
		return nil, err
	}

	return v, nil
}

func (_ *Int64Type) Write(_ godbase.Rec, _v interface{}, w io.Writer) error {
	v := _v.(int64)
	return WriteBinVal(8, &v, w)
}
